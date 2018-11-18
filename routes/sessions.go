package routes

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"hacking-portal/templates"

	"github.com/google/uuid"
	"gopkg.in/ldap.v2"
)

// Session stores user, token, and expiration
type Session struct {
	Username    string
	DisplayName string
	AccessLevel accessLevel
	Expires     time.Time
}

// Credentials stores username and password
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type accessLevel int

const (
	noAccess      accessLevel = 0
	studentAccess accessLevel = 1
	adminAccess   accessLevel = 2
)

var (
	sessions     map[string]Session
	adminUsers   map[string]bool
	ldapAddress  string
	ldapDC       string
	courseFilter string
)

func getLDAPConnection() (*ldap.Conn, error) {
	// return a secure LDAP connection
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapAddress, 389))
	if err != nil {
		return conn, err
	}

	return conn, conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
}

func authenticateUser(creds Credentials) error {
	// get an LDAP connection
	conn, err := getLDAPConnection()
	if err != nil {
		return err
	}

	// attempt to bind the credentials to LDAP, returning the possible error
	return conn.Bind(fmt.Sprintf("uid=%s,ou=people,%s", creds.Username, ldapDC), creds.Password)
}

func verifyAccess(username string) (accessLevel, error) {
	// return early if the user is a pre-defined admin
	if _, exists := adminUsers[username]; exists {
		return adminAccess, nil
	}

	// get an LDAP connection
	conn, err := getLDAPConnection()
	if err != nil {
		return noAccess, err
	}

	// get the course group information
	results, err := conn.Search(ldap.NewSearchRequest(
		"ou=groups,"+ldapDC,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		courseFilter, []string{}, nil,
	))
	if err != nil {
		return noAccess, err
	}

	// parse the results and find the user
	for _, entry := range results.Entries {
		for _, member := range entry.GetAttributeValues("memberUid") {
			if member == username {
				return studentAccess, nil
			}
		}
	}

	// user is not part of the course
	return noAccess, nil
}

func getUserDisplayName(username string) (string, error) {
	// get an LDAP connection
	conn, err := getLDAPConnection()
	if err != nil {
		return "", err
	}

	// search for user info
	results, err := conn.Search(ldap.NewSearchRequest(
		"ou=people,"+ldapDC,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username), []string{}, nil,
	))
	if err != nil {
		return "", err
	}

	// find and return the user's full name
	for _, entry := range results.Entries {
		name := entry.GetAttributeValue("displayName")
		if name != "" {
			return name, nil
		}
	}

	// fall back to the username
	return username, nil
}

// InitAuthentication sets required values for LDAP connection
func InitAuthentication(addr, dn, courseCode, admins string) {
	if addr == "" {
		log.Fatal("LDAP address must be provided")
	}
	if dn == "" {
		log.Fatal("LDAP domain name must be provided")
	}
	if courseCode == "" {
		log.Fatal("CourseCode must be provided")
	}
	if admins == "" {
		log.Fatal("Admin list must be provided")
	}

	adminUsers = make(map[string]bool)
	sessions = make(map[string]Session)

	for _, admin := range strings.Split(admins, " ") {
		adminUsers[admin] = true
	}

	ldapAddress = addr
	ldapDC = dn
	courseFilter = fmt.Sprintf("(cn=fs_%s_1)", courseCode)
}

// SessionHandler checks for a session and handles it accordingly
func SessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/admin") || strings.HasPrefix(r.URL.Path, "/group") {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				// no session
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			} else if _, exists := sessions[cookie.Value]; !exists {
				// we are not aware of the user's cookie, force em to log out
				http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
				return
			} else if sessions[cookie.Value].Expires.Before(time.Now()) {
				// session has expired, log the user out
				http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
				return
			}

			var path string
			session := sessions[cookie.Value]
			switch session.AccessLevel {
			case adminAccess:
				path = "/admin"
			case studentAccess:
				path = "/group"
			}

			if path != "" {
				// refresh the cookie
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   cookie.Value,
					Path:    "/",
					Expires: time.Now().Add(1 * time.Hour),
				})

				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
				} else {
					http.Redirect(w, r, path, http.StatusTemporaryRedirect)
				}
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// GetLogin routes invalid users to the login page
func GetLogin(w http.ResponseWriter, r *http.Request) {
	// prepare and ensure validity of template files
	tpl := template.Must(template.New("layout").Parse(templates.Layout + templates.Login))

	// render the templates
	tpl.ExecuteTemplate(w, "layout", nil)
}

// PostLogin validates the user and creates session
func PostLogin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := authenticateUser(creds); err != nil {
		log.Println(err)
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	level, err := verifyAccess(creds.Username)
	if err != nil || level == noAccess {
		http.Error(w, "inaccessible", http.StatusForbidden)
		return
	}

	name, err := getUserDisplayName(creds.Username)
	if err != nil {
		http.Error(w, "failed to get user name", http.StatusInternalServerError)
		return
	}

	sessionToken := uuid.New().String()
	expiration := time.Now().Add(1 * time.Hour)

	sessions[sessionToken] = Session{
		Username:    creds.Username,
		DisplayName: name,
		AccessLevel: level,
		Expires:     expiration,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: expiration,
	})

	fmt.Fprint(w, "OK")
}

// GetLogout handles user logouts
func GetLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// cookie already doesn't exist, just redirect
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// remove the session from our token storage
	delete(sessions, cookie.Value)

	// create a new, dead cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})

	// redirec to the login page
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// RedirectLogin sends the user to the login page
func RedirectLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

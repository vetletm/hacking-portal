package routes

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/ldap.v2"
)

// Session stores user, token, and expiration
type Session struct {
	Username string
	Name     string
	Status   SessionType
	Expiry   time.Time
}

// Credentials stores username and password
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// SessionType stores if user is student, admin, or invalid
type SessionType int

const (
	StudentUser SessionType = 1
	AdminUser   SessionType = 2
	InvalidUser SessionType = 0
)

// AuthenticationMiddleware keeps a map of authenticated users
type AuthenticationMiddleware struct {
	TokenUsers  map[string]Session
	Admins      map[string]bool
	LDAPAddress string
	LDAPDC      string
	GroupFilter string
}

func (amw *AuthenticationMiddleware) ValidateUser(creds Credentials) (SessionType, error) {
	// create a connection to LDAP
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", amw.LDAPAddress, 389))
	if err != nil {
		return InvalidUser, err
	}
	defer conn.Close()

	// switch the connection to TLS
	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return InvalidUser, err
	}

	// attempt to bind the credentials to LDAP
	username := fmt.Sprintf("uid=%s,ou=people,%s", creds.Username, amw.LDAPDC)
	err = conn.Bind(username, creds.Password)
	if err != nil {
		return InvalidUser, err
	}

	// return early if the user is an admin
	if _, exists := amw.Admins[creds.Username]; exists {
		return AdminUser, nil
	}

	// search for users in the course group
	groupResults, err := conn.Search(ldap.NewSearchRequest(
		"ou=groups,"+amw.LDAPDC,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		amw.GroupFilter,
		[]string{},
		nil,
	))
	if err != nil {
		return InvalidUser, err
	}

	// search for the user in the group
	for _, entry := range groupResults.Entries {
		for _, member := range entry.GetAttributeValues("memberUid") {
			if member == creds.Username {
				return StudentUser, nil
			}
		}
	}

	// user is not part of the group, bail out
	return InvalidUser, nil
}

func (amw *AuthenticationMiddleware) ValidateSession(token string) (SessionType, string) {
	session, found := amw.TokenUsers[token]
	if !found {
		return InvalidUser, ""
	}
	return session.Status, session.Username
}

var amw = AuthenticationMiddleware{}

func Init(addr, dn, courseCode, admins string) {
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

	amw.Admins = make(map[string]bool)
	amw.TokenUsers = make(map[string]Session)

	for _, admin := range strings.Split(admins, " ") {
		amw.Admins[admin] = true
	}

	amw.LDAPAddress = addr
	amw.LDAPDC = dn
	amw.GroupFilter = fmt.Sprintf("(cn=fs_%s_1)", courseCode)
}

func SessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/login") {
			next.ServeHTTP(w, r)
		} else {
			c, err := r.Cookie("session_token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			var path string
			sessionType, userID := amw.ValidateSession(c.Value)
			switch sessionType {
			case AdminUser:
				path = "/admin"
			case StudentUser:
				path = "/groups"
			}

			// inject userID into request context
			r = r.WithContext(context.WithValue(r.Context(), "session_user_id", userID))

			if path != "" {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
				} else {
					http.Redirect(w, r, path, http.StatusTemporaryRedirect)
				}
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
		}
	})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sessionType, err := amw.ValidateUser(creds)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	sessionToken := uuid.New().String()
	expiration := time.Now().Add(1 * time.Hour)

	amw.TokenUsers[sessionToken] = Session{
		Username: creds.Username,
		Status:   sessionType,
		Expiry:   expiration,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiration,
	})

	fmt.Fprint(w, "OK")
}

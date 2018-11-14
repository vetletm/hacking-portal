package routes

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SessionType int

// Session stores user, token, and expiration
type Session struct {
	UserName string
	Status   SessionType
	Expiry   string // nop
}

// Credentials stores username and password
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

const (
	StudentUser SessionType = 1
	AdminUser   SessionType = 2
	InvalidUser SessionType = 0
)

// AuthenticationMiddleware keeps a map of authenticated users
type AuthenticationMiddleware struct {
	TokenUsers map[string]Session
}

var amw = AuthenticationMiddleware{}

// func (amw *AuthenticationMiddleware) Populate() {
// 	amw.TokenUsers["00000000"] = models.User{"vetletm", "admin"}
// 	amw.TokenUsers["aaaaaaaa"] = models.User{"adrialu", "admin"}
// 	amw.TokenUsers["abababab"] = models.User{"miebkri", "student"}
// 	amw.TokenUsers["bcbcbcbc"] = models.User{"random", "student"}
// }
//
// func Init() {
// 	amw.Populate()
// }

func ValidateSession(token string) SessionType {
	session, found := amw.TokenUsers[token]
	if !found {
		return InvalidUser
	}
	return session.Status
}

func SessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session_token")
		// TODO: error handle cookie

		var path string
		switch ValidateSession(c.Value) {
		case AdminUser:
			path = "/admin"
		case StudentUser:
			path = "/student"
		}

		if path != "" {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}

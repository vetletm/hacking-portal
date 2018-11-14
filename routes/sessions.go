package routes

import (
	"hacking-portal/models"
	"log"
	"net/http"
)

// AuthenticationMiddleware keeps a map of authenticated users
type AuthenticationMiddleware struct {
	TokenUsers map[string]models.User
}

func (amw *AuthenticationMiddleware) Populate() {
	amw.TokenUsers["00000000"] = models.User{"vetletm", "admin"}
	amw.TokenUsers["aaaaaaaa"] = models.User{"adrialu", "admin"}
}

func Init() {
	amw := AuthenticationMiddleware{}
	amw.Populate()
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.TokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func GrabSession(username string) (string, bool) {

	return "vetletm", true
}

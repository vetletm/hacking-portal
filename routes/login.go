package routes

import "net/http"

//
func GetLogin(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - if the user doesn't have a session, show login form
	// - if the user has a session, redirect to the correct homepage
}

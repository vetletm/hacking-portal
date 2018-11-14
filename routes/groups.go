package routes

import (
	"net/http"

	"hacking-portal/db"

	"github.com/go-chi/chi"
)

// Groups is an implementation of the endpoint for all Groups-related methods.
// Database interfaces for all the methods are expected to be provided.
type GroupsEndpoint struct {
	Students db.StudentStorage
}

// GetGroups renders a view of all student groups
func (storage *GroupsEndpoint) GetGroups(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PostJoinGroup handles group join requests
func (storage *GroupsEndpoint) PostJoinGroup(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PostLeaveGroup handles group leave requests
func (storage *GroupsEndpoint) PostLeaveGroup(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GroupsRouter sets up routing for the group enrollment view
func GroupsRouter() chi.Router {
	ep := GroupsEndpoint{new(db.StudentDatabase)}

	r := chi.NewRouter()
	r.Get("/", ep.GetGroups)
	r.Post("/join", ep.PostJoinGroup)
	r.Post("/leave", ep.PostLeaveGroup)

	return r
}

package routes

import (
	"net/http"

	"hacking-portal/db"

	"github.com/go-chi/chi"
)

// GroupEndpoint is an implementation of the endpoint for all Group-related methods.
// Database interfaces for all the methods are expected to be provided.
type GroupEndpoint struct {
	Machines db.MachineStorage
	Students db.StudentStorage
}

// GetDashboard renders a view of the group interface
func (storage *GroupEndpoint) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GetMachineKey returns the given group's machine's PEM key
func (storage *GroupEndpoint) GetMachineKey(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PostMachineRestart handles a group's machine restart requests
func (storage *GroupEndpoint) PostMachineRestart(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GroupRouter sets up routing for the group dashboard view
func GroupRouter() chi.Router {
	ep := GroupEndpoint{
		Machines: new(db.MachineDatabase),
		Students: new(db.StudentDatabase),
	}

	r := chi.NewRouter()
	r.Get("/", ep.GetDashboard)
	r.Get("/key/{machineIndex:[0-9]+}", ep.GetMachineKey)
	r.Post("/restart/{machineIndex:[0-9]+}", ep.PostMachineRestart)

	return r
}

package routes

import (
	"net/http"

	"hacking-portal/db"

	"github.com/go-chi/chi"
)

// AdminEndpoint is an implementation of the endpoint for all Admin-related methods.
// Database interfaces for all the methods are expected to be provided.
type AdminEndpoint struct {
	Machines db.MachineStorage
	Students db.StudentStorage
}

// GetDashboard renders a view of the administration interface
func (storage *AdminEndpoint) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - shows buttons for pages
	//   - view machines
	//   - view tasks
	//   - view groups
}

// PostMachineAssign handles machine restart requests
func (storage *AdminEndpoint) PostMachineRestart(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - lists machines and their assigned groups
	//   - the assigned group is a dropdown that defaults to nothing
	//   - "assign" button next to the dropdown when a change is staged
	// - have some sorting? (machines/groups)
}

// PostMachineAssign handles machine group assignment requests
func (storage *AdminEndpoint) PostMachineAssign(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - assigns a machine to a group by ID
	// - if group ID is empty, unassign the machine from any group
	// - needs validation (just in case)

	// var groupID int
	// machineUUID := chi.URLParam(r, "machine")

	// if val := chi.URLParam(r, "id"); val != "" {
	// 	groupID, _ = strconv.Atoi(chi.URLParam(r, "id"))
	// }
}

// AdminRouter sets up routing for the administration web interface
func AdminRouter() chi.Router {
	ep := AdminEndpoint{
		Machines: new(db.MachineDatabase),
		Students: new(db.StudentDatabase),
	}

	r := chi.NewRouter()
	r.Get("/", ep.GetDashboard)
	r.Post("/restart/{machineUUID:[A-Za-z0-9-]+}", ep.PostMachineRestart)
	r.Post("/assign/{machineUUID:[A-Za-z0-9-]+}(?:/{groupID:[0-9]+})?", ep.PostMachineAssign)

	return r
}

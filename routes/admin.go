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

// GetHomepage renders a view for the administration web interface
func (storage *AdminEndpoint) GetHomepage(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - shows buttons for pages
	//   - view machines
	//   - view tasks
	//   - view groups
}

// GetMachines renders a view of all the machines in OpenStack
func (storage *AdminEndpoint) GetMachines(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - lists machines and their assigned groups
	//   - the assigned group is a dropdown that defaults to nothing
	//   - "assign" button next to the dropdown when a change is staged
	// - have some sorting? (machines/groups)
}

// PostAssign assigns a machine to a group
func (storage *AdminEndpoint) PostAssign(w http.ResponseWriter, r *http.Request) {
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
	r.Get("/", ep.GetHomepage)

	r.Route("/machines", func(r chi.Router) {
		r.Get("/", ep.GetMachines)
		r.Route("/{machine:[A-Za-z0-9-]+}", func(r chi.Router) {
			r.Post("/", ep.PostAssign)
			r.Post("/{id:[0-9]+}", ep.PostAssign)
		})
	})

	return r
}

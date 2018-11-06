package routes

import (
	"net/http"
	"strconv"

	"github.com/vetletm/hacking-portal/db"
)

// AdminEndpoint is an implementation of the endpoint for all Admin-related methods.
// Database interfaces for all the methods are expected to be provided.
type AdminEndpoint struct {
	answers  db.AnswerStorage
	machines db.MachineStorage
	students db.StudentStorage
	tasks    db.TaskStorage
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

	machineUUID := chi.URLParam(r, "machine")
	groupID := strconv.Atoi(chi.URLParam(r, "id"))
}

// GetTasks renders a view of all tasks
func (storage *AdminEndpoint) GetTasks(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - lists all tasks
	//   - with each task shows groups' status
	//     - green = completed AND correct
	//     - red = completed AND incorrect
	//     - white/gray = incomplete
	//   - maybe add statistics of how many groups have completed/successfully
	// - "new task" form
	//   - name
	//   - description
	//   - correct answer (for the list populating correct answers)
}

// NewTask creates a new task from form data
func (storage *AdminEndpoint) NewTask(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - takes a form to create a new task
	// - form validation
}

// EditTask modifies an existing task from form data
func (storage *AdminEndpoint) EditTask(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - changes an existing task
	// - if the description was changed, make sure students get a notification or smth next time they visit their homepage

	taskID := strconv.Atoi(chi.URLParam(r, "task"))
}

// GetGroups renders a view of all groups
func (storage *AdminEndpoint) GetGroups(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - lists all groups
}

// GetGroup renders a view of a single group
func (storage *AdminEndpoint) GetGroup(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - shows single group information
	//   - group name
	//   - group members (name + studentID + username)
	//   - which machines they have access to
	//   - tasks status (like the tasks list)

	groupID := strconv.Atoi(chi.URLParam(r, "id"))
}

// AdminRouter sets up routing for the administration web interface
func AdminRouter() chi.Router {
	ep := AdminEndpoint{new(db.AdminDatabase)}

	r := chi.NewRouter()
	r.Get("/", ep.GetHomepage)

	r.Route("/machines", func(r chi.Router) {
		r.Get("/", ep.GetMachines)
		r.Post("/{machine:[A-Za-z0-9-]+}", ep.PostUnassign)
		r.Post("/{machine:[A-Za-z0-9-]+}/{id:[0-9]+}", ep.PostAssign)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", ep.GetTasks)
		r.Post("/", ep.NewTask)
		r.Put("/{id:[0-9]+}", ep.EditTask)
	})

	r.Route("/groups", func(r chi.Router) {
		r.Get("/", ep.GetGroups)
		r.Get("/{task:[0-9]+}", ep.GetGroup)
	})

	return r
}

package routes

import (
	"html/template"
	"net/http"
	"path"

	"hacking-portal/db"
	"hacking-portal/models"

	"github.com/go-chi/chi"
)

// GroupEndpoint is an implementation of the endpoint for all Group-related methods.
// Database interfaces for all the methods are expected to be provided.
type GroupEndpoint struct {
	Machines db.MachineStorage
	Students db.StudentStorage
}

type groupPageData struct {
	User     models.Student
	Machines []models.Machine
}

// GetDashboard renders a view of the group interface
func (storage *GroupEndpoint) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// get the user from the session (type-casted)
	username := r.Context().Value("session_user_id").(string)

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(username)
	if err != nil || sessionUser.GroupID == 0 {
		// sessionUser doesn't exist or has no group affiliation, redirect em'
		http.Redirect(w, r, "/groups", http.StatusTemporaryRedirect)
		return
	}

	// prepare page data
	pageData := groupPageData{User: sessionUser}

	// get the machines
	// TODO: get from OpenStack
	pageData.Machines = []models.Machine{
		{"123", 0, 1, "10.212.136.10"},
		{"456", 0, 2, "10.212.136.20"},
		{"789", 0, 3, "10.212.136.30"},
	}

	// prepare and ensure validity of template files
	tpl := template.Must(template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "navigation.html"),
		path.Join("templates", "group.html"),
	))

	// render the templates with data
	tpl.ExecuteTemplate(w, "layout", pageData)
}

// GetMachineKey returns the given group's machine's PEM key
func (storage *GroupEndpoint) GetMachineKey(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PostMachineRestart handles a group's machine restart requests
func (storage *GroupEndpoint) PostMachineRestart(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PostLeaveGroup handles group leave requests
func (storage *GroupEndpoint) GetLeaveGroup(w http.ResponseWriter, r *http.Request) {
	// get the user from the session (type-casted)
	username := r.Context().Value("session_user_id").(string)

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(username)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusBadRequest)
		return
	}

	// attempt to get the student information, validating it
	if student, err := storage.Students.FindByID(sessionUser.ID); err != nil {
		http.Error(w, "Unable to get student data", http.StatusInternalServerError)
	} else if student.GroupID == 0 {
		http.Error(w, "Student is not in a group", http.StatusBadRequest)
	} else {
		student.GroupID = 0

		if err := storage.Students.Upsert(student); err != nil {
			http.Error(w, "Unable to leave group", http.StatusInternalServerError)
		} else {
			// redirect to the groups view
			http.Redirect(w, r, "/groups", http.StatusTemporaryRedirect)
		}
	}
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
	r.Get("/leave", ep.GetLeaveGroup)

	return r
}

package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"hacking-portal/db"
	"hacking-portal/models"
	"hacking-portal/openstack"

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
	if pageData.Machines, err = storage.Machines.FindByGroup(sessionUser.GroupID); err != nil {
		http.Error(w, "unable to get machines", http.StatusInternalServerError)
		return
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
	// get uuid from URL path
	uuid := chi.URLParam(r, "machineUUID")

	// get the user from the session (type-casted)
	username := r.Context().Value("session_user_id").(string)

	// Compare requested machine's group id to user's group id, and reboot
	if sessionUser, err := storage.Students.FindByID(username); err != nil {
		http.Error(w, "Invalid user session", http.StatusBadRequest)
	} else if sessionUser.GroupID == 0 {
		http.Error(w, "Invalid user session", http.StatusBadRequest)
	} else if machine, err := storage.Machines.FindByID(uuid); err != nil {
		http.Error(w, "Invalid machine", http.StatusBadRequest)
	} else if machine.GroupID != sessionUser.GroupID {
		http.Error(w, "Invalid machine", http.StatusBadRequest)
	} else if openstack.Reboot(uuid) != nil {
		http.Error(w, "Could not reboot machine", http.StatusInternalServerError)
	} else {
		fmt.Fprint(w, "OK")
	}
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
	r.Get("/key/{machineUUID:[A-Za-z0-9-]+}", ep.GetMachineKey)
	r.Post("/restart/{machineUUID:[A-Za-z0-9-]+}", ep.PostMachineRestart)
	r.Get("/leave", ep.GetLeaveGroup)

	return r
}

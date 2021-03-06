package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"hacking-portal/db"
	"hacking-portal/models"
	"hacking-portal/openstack"
	"hacking-portal/templates"

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
	// get the user from the session
	cookie, _ := r.Cookie("session_token")
	session := sessions[cookie.Value]

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(session.Username)
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

	// add some handy functions to the template engine
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	// prepare and ensure validity of template files
	tpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.Layout + templates.Navigation + templates.Group))

	// render the templates with data
	tpl.ExecuteTemplate(w, "layout", pageData)
}

// PostMachineRestart handles a group's machine restart requests
func (storage *GroupEndpoint) PostMachineRestart(w http.ResponseWriter, r *http.Request) {
	// get the user from the session
	cookie, _ := r.Cookie("session_token")
	session := sessions[cookie.Value]

	// get uuid from URL path
	uuid := chi.URLParam(r, "machineUUID")

	// Compare requested machine's group id to user's group id, and reboot
	if sessionUser, err := storage.Students.FindByID(session.Username); err != nil {
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

// GetLeaveGroup handles group leave requests
func (storage *GroupEndpoint) GetLeaveGroup(w http.ResponseWriter, r *http.Request) {
	// get the user from the session
	cookie, _ := r.Cookie("session_token")
	session := sessions[cookie.Value]

	// attempt to get the student information, validating it
	if student, err := storage.Students.FindByID(session.Username); err != nil {
		http.Error(w, "Unable to get student data", http.StatusBadRequest)
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
	r.Post("/restart/{machineUUID:[A-Za-z0-9-]+}", ep.PostMachineRestart)
	r.Get("/leave", ep.GetLeaveGroup)

	return r
}

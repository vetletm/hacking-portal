package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"sort"

	"hacking-portal/db"
	"hacking-portal/models"
	"hacking-portal/openstack"

	"github.com/go-chi/chi"
)

// AdminEndpoint is an implementation of the endpoint for all Admin-related methods.
// Database interfaces for all the methods are expected to be provided.
type AdminEndpoint struct {
	Machines db.MachineStorage
	Students db.StudentStorage
}

type adminPageData struct {
	User     models.Student
	Machines []models.Machine
	Groups   []models.Group
}

// GetDashboard renders a view of the administration interface
func (storage *AdminEndpoint) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// get the user from the session (type-casted)
	username := r.Context().Value("session_user_id").(string)

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(username)
	if err != nil {
		// sessionUser doesn't exist yet, we'll have to create it
		// this will happen on first visit
		sessionUser = models.Student{ID: username}

		err = storage.Students.Upsert(sessionUser)
		if err != nil {
			// something went horribly wrong
			http.Error(w, "Unable to initiate user", http.StatusInternalServerError)
			return
		}
	}

	// prepare page data
	pageData := adminPageData{User: sessionUser}

	// get the groups
	if groups, err := storage.Students.FindGroups(); err != nil {
		http.Error(w, "Unable to get groups", http.StatusInternalServerError)
		return
	} else {
		// maps are intentionally randomized in order, so we have to get an ordered slice of it
		var groupKeys []int
		for key := range groups {
			groupKeys = append(groupKeys, key)
		}
		sort.Ints(groupKeys)

		// iterate over each group and fill in the page data
		for _, groupID := range groupKeys {
			// append the group data and members to the page data
			pageData.Groups = append(pageData.Groups, models.Group{ID: groupID})
		}
	}

	// get the machines from the database
	if pageData.Machines, err = storage.Machines.FindAll(); err != nil {
		http.Error(w, "unable to grab machines", http.StatusInternalServerError)
		return
	}

	// prepare and ensure validity of template files
	tpl := template.Must(template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "navigation.html"),
		path.Join("templates", "admin.html"),
	))

	// render the templates with data
	tpl.ExecuteTemplate(w, "layout", pageData)
}

// PostMachineAssign handles machine restart requests
func (storage *AdminEndpoint) PostMachineRestart(w http.ResponseWriter, r *http.Request) {
	// get machine UUID from URL path
	uuid := chi.URLParam(r, "machineUUID")

	// - lists machines and their assigned groups
	if _, err := storage.Machines.FindByID(uuid); err != nil {
		http.Error(w, "Couldn't get machines from db", http.StatusNotFound)
		return
	}

	// Attempt to reboot the server
	if openstack.Reboot(uuid) != nil {
		http.Error(w, "Failed to reboot machine", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}

// PostMachineAssign handles machine group assignment requests
func (storage *AdminEndpoint) PostMachineAssign(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	// attempt to decode and validate the body contents, then get the machine information
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if groupID, ok := payload["groupID"]; !ok {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if machineUUID, ok := payload["machineUUID"]; !ok {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if machine, err := storage.Machines.FindByID(machineUUID.(string)); err != nil {
		http.Error(w, "Could not find machine", http.StatusNotFound)
	} else {
		// Set new group id
		machine.GroupID = groupID.(int)

		// attempt to update the machine in database
		if storage.Machines.Upsert(machine) != nil {
			http.Error(w, "Could not update machine", http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, "OK")
		}
	}
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
	r.Post("/assign", ep.PostMachineAssign)

	return r
}

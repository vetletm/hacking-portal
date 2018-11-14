package routes

import (
	"html/template"
	"net/http"
	"path"
	"sort"

	"hacking-portal/db"
	"hacking-portal/models"

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

	// get the machines
	// TODO: get from OpenStack
	pageData.Machines = []models.Machine{
		{"123", 0, 1, "10.212.136.10"},
		{"456", 0, 2, "10.212.136.20"},
		{"789", 1, 1, "10.212.136.30"},
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

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

// Groups is an implementation of the endpoint for all Groups-related methods.
// Database interfaces for all the methods are expected to be provided.
type GroupsEndpoint struct {
	Students db.StudentStorage
}

type groupsPageData struct {
	User   models.Student
	Groups []models.Group
}

// GetGroups renders a view of all student groups
func (storage *GroupsEndpoint) GetGroups(w http.ResponseWriter, r *http.Request) {
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

	if sessionUser.GroupID != 0 {
		// the user already has a group, redirect em
		http.Redirect(w, r, "/group", http.StatusTemporaryRedirect)
		return
	}

	// prepare page data
	pageData := groupsPageData{User: sessionUser}

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
			numMembers := groups[groupID]

			// get all group members
			if groupMembers, err := storage.Students.FindByGroup(groupID); err != nil {
				http.Error(w, "Unable to parse groups", http.StatusInternalServerError)
			} else {
				// append the group data and members to the page data
				pageData.Groups = append(pageData.Groups, models.Group{
					ID:      groupID,
					Full:    numMembers == 3, // hardcode much
					Members: groupMembers,
				})
			}
		}
	}

	// append empty group at the end so people can join an empty one
	nextGroupID := len(pageData.Groups) + 1
	pageData.Groups = append(pageData.Groups, models.Group{
		ID: nextGroupID,
	})

	// prepare and ensure validity of template files
	tpl := template.Must(template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "navigation.html"),
		path.Join("templates", "groups.html"),
	))

	// render the templates with data
	tpl.ExecuteTemplate(w, "layout", pageData)
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

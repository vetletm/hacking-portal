package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"sort"

	"hacking-portal/db"
	"hacking-portal/models"
	"hacking-portal/templates"

	"github.com/go-chi/chi"
)

// GroupsEndpoint interfaces for all the methods are expected to be provided.
type GroupsEndpoint struct {
	Students db.StudentStorage
}

type groupsPageData struct {
	User   models.Student
	Groups []models.Group
}

// GetGroups renders a view of all student groups
func (storage *GroupsEndpoint) GetGroups(w http.ResponseWriter, r *http.Request) {
	// get the user from the session
	cookie, _ := r.Cookie("session_token")
	session := sessions[cookie.Value]

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(session.Username)
	if err != nil {
		// sessionUser doesn't exist yet, we'll have to create it
		// this will happen on first visit
		sessionUser = models.Student{ID: session.Username, Name: session.DisplayName}

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
	groups, err := storage.Students.FindGroups()
	if err != nil {
		http.Error(w, "Unable to get groups", http.StatusInternalServerError)
		return
	}

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

	// append empty group at the end so people can join an empty one
	nextGroupID := len(pageData.Groups) + 1
	pageData.Groups = append(pageData.Groups, models.Group{
		ID: nextGroupID,
	})

	// prepare and ensure validity of template files
	tpl := template.Must(template.New("layout").Parse(templates.Layout + templates.Navigation + templates.Groups))

	// render the templates with data
	tpl.ExecuteTemplate(w, "layout", pageData)
}

// PostJoinGroup handles group join requests
func (storage *GroupsEndpoint) PostJoinGroup(w http.ResponseWriter, r *http.Request) {
	// get the user from the session
	cookie, _ := r.Cookie("session_token")
	session := sessions[cookie.Value]

	// get the actual sessionUser object from the username
	sessionUser, err := storage.Students.FindByID(session.Username)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusBadRequest)
		return
	}

	if sessionUser.GroupID != 0 {
		http.Error(w, "User already in a group", http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}

	// attempt to decode and validate the body contents, then get the student information
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if groupID, ok := payload["groupID"]; !ok {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if reflect.TypeOf(groupID).Kind() != reflect.Float64 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
	} else if groupID.(float64) <= 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
	} else {
		sessionUser.GroupID = int(groupID.(float64))

		// attempt to update the sessionUser's group ID
		if err := storage.Students.Upsert(sessionUser); err != nil {
			http.Error(w, "Unable to join group", http.StatusInternalServerError)
		} else {
			// render a successful message
			fmt.Fprint(w, "OK")
		}
	}
}

// GroupsRouter sets up routing for the group enrollment view
func GroupsRouter() chi.Router {
	ep := GroupsEndpoint{new(db.StudentDatabase)}

	r := chi.NewRouter()
	r.Get("/", ep.GetGroups)
	r.Post("/join", ep.PostJoinGroup)

	return r
}

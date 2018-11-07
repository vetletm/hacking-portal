package routes

import (
	"net/http"
	"strconv"

	"hacking-portal/db"

	"github.com/go-chi/chi"
)

// StudentEndpoint is an implementation of the endpoint for all Student-related methods.
// Database interfaces for all the methods are expected to be provided.
type StudentEndpoint struct {
	Answers  db.AnswerStorage
	Machines db.MachineStorage
	Students db.StudentStorage
	Tasks    db.TaskStorage
}

// GetHomepage renders a view for the student web interface
func (storage *StudentEndpoint) GetHomepage(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - if not in a group, redirect to /student/groups
	// - if in a group, list tasks and machines
	//   - tasks should show description + answer field + submit button
	//     - if a task description was changed since last time (not from POST but PUT) show a notification or smth
	//   - machines should show restart button + name + IP + download button for the PEM file
}

// GetGroups renders a view of all groups
func (storage *StudentEndpoint) GetGroups(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - list available groups and a "create new group" form
	// - should not be available to students already in a group
	//   - (redirect to /student)
}

// PostGroupJoin assigns the authenticated student to a group
func (storage *StudentEndpoint) PostGroupJoin(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - assigns the current student with the given group
	// - should only work with students that are not in a group
	//   - (respond somehow nicely)

	groupID := strconv.Atoi(chi.URLParam(r, "id"))
}

// PostGroupCreate creates a new group and assigns the authenticated student to it
func (storage *StudentEndpoint) PostGroupCreate(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - takes form data and creates a group AND makes the current student join said group
	// - need to validate form
	// - should only work with students that are not in a group
	//   - (respond somehow nicely)

	groupID := strconv.Atoi(chi.URLParam(r, "id"))
}

// PostGroupLeave removes the authenticated student from its current group
func (storage *StudentEndpoint) PostGroupLeave(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - removes a student from a group
	// - should only work with students that are in a group
	//   - (respond somehow nicely)
}

// PostRestart restarts a given OpenStack instance by UUID
func (storage *StudentEndpoint) PostRestart(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - restarts the OpenStack instance with given id
	//   - (respond somehow nicely)
	// - should only work for IDs the group has access to
	//   - (respond somehow nicely)

	machineUUID := chi.URLParam(r, "machine")
}

// PostTask stores a task's answer for the authenticated student's group
func (storage *StudentEndpoint) PostTask(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// - lets a group submit a task result

	taskID := strconv.Atoi(chi.URLParam(r, "task"))
}

// StudentRouter sets up routing for the student web interface
func StudentRouter() chi.Router {
	ep := StudentEndpoint{
		Answers:  new(db.AnswerDatabase),
		Machines: new(db.MachineDatabase),
		Students: new(db.StudentDatabase),
		Tasks:    new(db.TaskDatabase),
	}

	r := chi.NewRouter()
	r.Get("/", ep.GetHomepage)
	r.Route("/groups", func(r chi.Router) {
		r.Post("/join/{id:[0-9]+}", ep.PostGroupJoin)
		r.Post("/create/{id:[0-9]+}", ep.PostGroupCreate)
		r.Post("/leave", ep.PostGroupLeave)
	})

	r.Post("/restart/{machine:[A-Za-z0-9-]+}", ep.PostRestart)
	r.Post("/task/{task:[0-9]+}", ep.PostTask)

	return r
}

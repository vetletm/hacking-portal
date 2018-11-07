package db

import (
	"hacking-portal/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// tasks collection structure:
// { id, description, answer, deadline }

// TaskStorage is an interface describing the methods of the TaskDatabase struct
type TaskStorage interface {
	FindByID(int) (models.Task, error)
	FindTasks() ([]models.Task, error)
	Upsert(models.Task) (models.Task, error)
}

// TaskDatabase is an implementation of the storage for all Task-related methods
type TaskDatabase struct{}

// FindByID returns a single task by ID
func (TaskDatabase) FindByID(id int) (models.Task, error) {
	var task models.Task
	err := db.C("tasks").Find(bson.M{"id": id}).One(&task)
	return task, err
}

// FindGroups returns an array of all the group IDs
func (TaskDatabase) FindTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := db.C("tasks").Find(nil).All(&tasks)
	return tasks, err
}

// Upsert adds/updates the task to the database and returns the updated version
func (TaskDatabase) Upsert(task models.Task) (models.Task, error) {
	// insert or update the task
	_, err := db.C("tasks").Find(bson.M{"id": task.ID}).Apply(mgo.Change{
		Update:    task,
		Upsert:    true,
		ReturnNew: true,
	}, &task)

	return task, err
}

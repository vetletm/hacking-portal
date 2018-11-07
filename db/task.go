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
	FindAll() ([]models.Task, error)
	FindByID(int) (models.Task, error)
	Upsert(models.Task) (models.Task, error)
}

// TaskDatabase is an implementation of the storage for all Task-related methods
type TaskDatabase struct{}

// FindAll returns an array of all the group IDs
func (TaskDatabase) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := db.C("tasks").Find(nil).All(&tasks)
	return tasks, err
}

// FindByID returns a single task by ID
func (TaskDatabase) FindByID(id int) (models.Task, error) {
	var task models.Task
	err := db.C("tasks").Find(bson.M{"id": id}).One(&task)
	return task, err
}

// Upsert adds/updates the task to the database and returns the updated version
func (TaskDatabase) Upsert(task models.Task) (models.Task, error) {
	// get the next available task ID for when creating a new task
	if task.ID == 0 {
		var ids map[string]int
		if _, err := db.C("counters").Find(nil).Apply(mgo.Change{
			Update:    bson.M{"$inc": bson.M{"nextTaskID": 1}},
			Upsert:    true,
			ReturnNew: true,
		}, &ids); err != nil {
			return task, err
		}

		task.ID = ids["nextTaskID"]
	}

	// insert or update the task
	_, err := db.C("tasks").Find(bson.M{"id": task.ID}).Apply(mgo.Change{
		Update:    task,
		Upsert:    true,
		ReturnNew: true,
	}, &task)

	return task, err
}

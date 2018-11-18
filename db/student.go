package db

import (
	"hacking-portal/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// students collection structure:
// { id, alias, name, groupID }

// StudentStorage is an interface describing the methods of the StudentDatabase struct
type StudentStorage interface {
	FindAll() ([]models.Student, error)
	FindByID(string) (models.Student, error)
	FindByName(string) (models.Student, error)
	FindByGroup(int) ([]models.Student, error)
	FindGroups() (map[int]int, error)
	Upsert(models.Student) error
}

// StudentDatabase is an implementation of the storage for all Student-related methods
type StudentDatabase struct{}

// FindAll returns an array of all the students
func (StudentDatabase) FindAll() ([]models.Student, error) {
	var students []models.Student
	err := db.C("students").Find(nil).All(&students)
	return students, err
}

// FindByID returns a single student by ID (username)
func (StudentDatabase) FindByID(id string) (models.Student, error) {
	var student models.Student
	err := db.C("students").Find(bson.M{"id": id}).One(&student)
	return student, err
}

// FindByName returns a single student by name
func (StudentDatabase) FindByName(name string) (models.Student, error) {
	var student models.Student
	err := db.C("students").Find(bson.M{"name": name}).One(&student)
	return student, err
}

// FindByGroup finds all students in a certain group
func (StudentDatabase) FindByGroup(groupID int) ([]models.Student, error) {
	var students []models.Student
	err := db.C("students").Find(bson.M{"groupID": groupID}).All(&students)
	return students, err
}

// FindGroups returns a map of all group IDs with the number of members
func (StudentDatabase) FindGroups() (map[int]int, error) {
	groupIDs := map[int]int{}
	var students []models.Student

	// get all students from the database
	if err := db.C("students").Find(nil).All(&students); err != nil {
		return groupIDs, err
	}

	// populate array with unique group IDs
	for _, student := range students {
		groupID := student.GroupID
		if groupID != 0 {
			if _, isset := groupIDs[groupID]; !isset {
				groupIDs[groupID] = 1
			} else {
				groupIDs[groupID]++
			}
		}
	}

	return groupIDs, nil
}

// Upsert adds/updates the student to the database
func (StudentDatabase) Upsert(student models.Student) error {
	_, err := db.C("students").Find(bson.M{"id": student.ID}).Apply(mgo.Change{
		Update: student,
		Upsert: true,
	}, nil)

	return err
}

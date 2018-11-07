package db

import (
	"sort"

	"hacking-portal/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// students collection structure:
// { id, alias, name, groupID }

// StudentStorage is an interface describing the methods of the StudentDatabase struct
type StudentStorage interface {
	FindAll() ([]models.Student, error)
	FindByID(int) (models.Student, error)
	FindByAlias(string) (models.Student, error)
	FindByName(string) (models.Student, error)
	FindByGroup(int) ([]models.Student, error)
	FindGroups() ([]int, error)
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

// FindByID returns a single student by ID
func (StudentDatabase) FindByID(id int) (models.Student, error) {
	var student models.Student
	err := db.C("students").Find(bson.M{"id": id}).One(&student)
	return student, err
}

// FindByAlias returns a single student by alias (username)
func (StudentDatabase) FindByAlias(alias string) (models.Student, error) {
	var student models.Student
	err := db.C("students").Find(bson.M{"alias": alias}).One(&student)
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

// FindGroups returns an array of all group IDs
func (StudentDatabase) FindGroups() ([]int, error) {
	var groupIDs []int
	var students []models.Student

	// get all students from the database
	if err := db.C("students").Find(nil).All(&students); err != nil {
		return groupIDs, err
	}

	// populate array with unique group IDs
	var groupExists []bool
	for _, student := range students {
		if student.GroupID != 0 && !groupExists[student.GroupID] {
			groupIDs = append(groupIDs, student.GroupID)
			groupExists[student.GroupID] = true
		}
	}

	// sort the array before returning
	sort.Ints(groupIDs)
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

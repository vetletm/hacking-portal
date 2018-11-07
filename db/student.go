package db

import (
	"hacking-portal/models"

	"github.com/globalsign/mgo/bson"
)

// StudentStorage is an interface describing the methods of the StudentDatabase struct
type StudentStorage interface {
	FindAll() ([]models.Student, error)
	FindByID(int) (models.Student, error)
	FindByAlias(string) (models.Student, error)
	FindByName(string) (models.Student, error)
	FindByGroup(int) ([]models.Student, error)
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

// Upsert adds/updates the student to the database
func (StudentDatabase) Upsert(student models.Student) error {
	// TODO
	return nil
}

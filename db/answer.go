package db

import (
	"hacking-portal/models"

	"github.com/globalsign/mgo/bson"
)

// AnswerStorage is an interface describing the methods of the AnswerDatabase struct
type AnswerStorage interface {
	FindAll() ([]models.Answer, error)
	FindByID(int) (models.Answer, error)
	FindByGroup(int) ([]models.Answer, error)
	Upsert(models.Answer) error
}

// AnswerDatabase is an implementation of the storage for all Answer-related methods
type AnswerDatabase struct{}

// FindAll returns an array of all the answers
func (AnswerDatabase) FindAll() ([]models.Answer, error) {
	var answers []models.Answer
	err := db.C("answers").Find(nil).All(&answers)
	return answers, err
}

// FindByID returns a single answer by ID
func (AnswerDatabase) FindByID(id int) (models.Answer, error) {
	var answer models.Answer
	err := db.C("answers").Find(bson.M{"taskID": id}).One(&answer)
	return answer, err
}

// FindByGroup finds all answers in a certain group
func (AnswerDatabase) FindByGroup(groupID int) ([]models.Answer, error) {
	var answers []models.Answer
	err := db.C("answers").Find(bson.M{"groupID": groupID}).All(&answers)
	return answers, err
}

// Upsert adds/updates the answer to the database
func (AnswerDatabase) Upsert(answer models.Answer) error {
	// TODO
	return nil
}

package models

// Student stores information about the student
type Student struct {
	ID      string `bson:"id"`
	Name    string `bson:"name"`
	GroupID int    `bson:"groupID"`
}

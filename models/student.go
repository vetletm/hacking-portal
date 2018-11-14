package models

type Student struct {
	ID      string `bson:"id"`
	Name    string `bson:"name"`
	GroupID int    `bson:"groupID"`
}

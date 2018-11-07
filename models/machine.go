package models

type Machine struct {
	ID      string `bson:"uuid"`
	Name    string `bson:"name"`
	GroupID int    `bson:"groupID"`
}

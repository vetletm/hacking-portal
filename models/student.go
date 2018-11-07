package models

type Student struct {
	ID      int    `bson:"id"`
	GroupID int    `bson:"groupID"`
	Name    string `bson:"name"`
	Alias   string `bson:"alias"`
}

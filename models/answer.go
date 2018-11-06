package models

type Answer struct {
	TaskID  int    `bson:"taskID"`
	GroupID int    `bson:"groupID"`
	Answer  string `bson:"answer"`
}

package models

type Student struct {
	GroupID   int    `bson:"groupID"`
	StudentID int    `bson:"studentID"`
	Name      string `bson:"studentName"`
	Alias     string `bson:"studentAlias"`
}

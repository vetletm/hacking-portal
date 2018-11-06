package models

type Task struct {
	ID          int    `bson:"taskID"`
	Description string `bson:"taskDescription"`
	Answer      string `bson:"answer"`
	Deadline    date   `bson:"deadline"`
}

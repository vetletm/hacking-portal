package models

import "time"

type Task struct {
	ID          int       `bson:"taskID"`
	Description string    `bson:"taskDescription"`
	Answer      string    `bson:"answer"`
	Deadline    time.Time `bson:"deadline"`
}

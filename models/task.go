package models

import "time"

type Task struct {
	ID          int       `bson:"id"`
	Description string    `bson:"description"`
	Answer      string    `bson:"answer"`
	Deadline    time.Time `bson:"deadline"`
}

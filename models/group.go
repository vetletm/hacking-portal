package models

type Group struct {
	ID      int
	Full    bool
	Members []Student
}

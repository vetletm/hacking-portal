package models

// Group stores information about a given group
type Group struct {
	ID      int
	Full    bool
	Members []Student
}

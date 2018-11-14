package models

// User stores the username and status of logged in user
type User struct {
	Username string // Username used to authenticate to LDAP
	Status   string // Either student or administrator
}

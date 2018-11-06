package models

type Machine struct {
	ID      string `bson:"machineID"`
	Name    string `bson:"machineName"`
	GroupID int    `bson:"groupID"`
	PEM     string `bson:"PEM"`
}

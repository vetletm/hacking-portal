package models

// Machine stores information about server in OpenStack
type Machine struct {
	Name       string `bson:"name"`
	UUID       string `bson:"uuid"`
	GroupID    int    `bson:"groupID"`
	GroupIndex int    `bson:"groupIndex"`
	Address    string `bson:"address"`
}

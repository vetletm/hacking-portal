package models

type Machine struct {
	UUID       string `bson:"uuid"`
	GroupID    int    `bson:"groupID"`
	GroupIndex int    `bson:"groupIndex"`
	Address    string
}

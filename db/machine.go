package db

import (
	"hacking-portal/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// machines collection structure:
// { uuid, name, groupID }

// MachineStorage is an interface describing the methods of the MachineDatabase struct
type MachineStorage interface {
	FindAll() ([]models.Machine, error)
	FindByID(string) (models.Machine, error)
	FindByGroup(int) ([]models.Machine, error)
	Upsert(models.Machine) error
}

// MachineDatabase is an implementation of the storage for all Machine-related methods
type MachineDatabase struct{}

// FindAll returns an array of all the machines
func (MachineDatabase) FindAll() ([]models.Machine, error) {
	var machines []models.Machine
	err := db.C("machines").Find(nil).All(&machines)
	return machines, err
}

// FindByID returns a single machine by ID
func (MachineDatabase) FindByID(uuid string) (models.Machine, error) {
	var machine models.Machine
	err := db.C("machines").Find(bson.M{"uuid": uuid}).One(&machine)
	return machine, err
}

// FindByGroup finds all machines in a certain group
func (MachineDatabase) FindByGroup(groupID int) ([]models.Machine, error) {
	var machines []models.Machine
	err := db.C("machines").Find(bson.M{"groupID": groupID}).All(&machines)
	return machines, err
}

// Upsert adds/updates the machine to the database
func (MachineDatabase) Upsert(machine models.Machine) error {
	_, err := db.C("machines").Find(bson.M{"uuid": machine.ID}).Apply(mgo.Change{
		Update: machine,
		Upsert: true,
	}, nil)

	return err
}

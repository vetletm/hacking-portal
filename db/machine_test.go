package db

import (
	"testing"

	"hacking-portal/models"

	"github.com/stretchr/testify/require"
)

func TestMachineUpsert(t *testing.T) {
	// new database type
	tdb := MachineDatabase{}

	// attempt to insert
	err := tdb.Upsert(models.Machine{
		UUID:    "1234",
		Name:    "foo",
		GroupID: 1,
	})

	// assert output
	require.Nil(t, err, "failed to insert machine")
}

func TestMachineFindAll(t *testing.T) {
	// new database type
	tdb := MachineDatabase{}

	// attempt to find all machines
	machines, err := tdb.FindAll()

	// assert output
	require.Nil(t, err, "failed to get machines")
	require.Len(t, machines, 1) // this runs after upsert, so there should be 1
}

func TestMachineFindByID(t *testing.T) {
	// new database type
	tdb := MachineDatabase{}

	// attempt to find machine by ID
	machine, err := tdb.FindByID("1234") // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get single machine")
	require.EqualValues(t, "foo", machine.Name)
}

func TestMachineFindByGroup(t *testing.T) {
	// new database type
	tdb := MachineDatabase{}

	// attempt to find machines by group
	machines, err := tdb.FindByGroup(1) // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get machines")
	require.Len(t, machines, 1) // this runs after upsert, so there should be 1
}

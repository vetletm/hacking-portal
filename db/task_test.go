package db

import (
	"testing"

	"hacking-portal/models"

	"github.com/stretchr/testify/require"
)

func TestTaskUpsert(t *testing.T) {
	// new database type
	tdb := TaskDatabase{}

	// attempt to insert
	task, err := tdb.Upsert(models.Task{
		Description: "foo",
		Answer:      "bar",
	})

	// assert output
	require.Nil(t, err, "failed to insert task")
	require.EqualValues(t, 1, task.ID)
	require.EqualValues(t, "foo", task.Description)
	require.EqualValues(t, "bar", task.Answer)
}

func TestTaskFindAll(t *testing.T) {
	// new database type
	tdb := TaskDatabase{}

	// attempt to find tasks
	tasks, err := tdb.FindAll()

	// assert output
	require.Nil(t, err, "failed to get tasks")
	require.Len(t, tasks, 1) // this runs after upsert, so there should be 1
}

func TestTaskFindByID(t *testing.T) {
	// new database type
	tdb := TaskDatabase{}

	// attempt to find tasks
	task, err := tdb.FindByID(1)

	// assert output
	require.Nil(t, err, "failed to get task")
	require.EqualValues(t, "foo", task.Description)
	require.EqualValues(t, "bar", task.Answer)
}

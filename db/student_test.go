package db

import (
	"testing"

	"hacking-portal/models"

	"github.com/stretchr/testify/require"
)

func TestStudentUpsert(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	var err error
	// attempt to insert a few students, asserting the output
	err = tdb.Upsert(models.Student{1, 10, "One One", "one"})
	require.Nil(t, err, "failed to insert student 1")
	err = tdb.Upsert(models.Student{2, 0, "Two Two", "two"})
	require.Nil(t, err, "failed to insert student 2")
	err = tdb.Upsert(models.Student{3, 10, "Three Three", "three"})
	require.Nil(t, err, "failed to insert student 3")
}

func TestStudentFindAll(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find students
	students, err := tdb.FindAll()

	// assert output
	require.Nil(t, err, "failed to get students")
	require.Len(t, students, 3) // this runs after upsert, so there should be 3
}

func TestStudentFindByID(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find student by ID
	student, err := tdb.FindByID(1) // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get student")
	require.EqualValues(t, 1, student.ID)
	require.EqualValues(t, 10, student.GroupID)
	require.EqualValues(t, "One One", student.Name)
	require.EqualValues(t, "one", student.Alias)
}

func TestStudentFindByAlias(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find student by alias
	student, err := tdb.FindByAlias("two") // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get student")
	require.EqualValues(t, 2, student.ID)
	require.EqualValues(t, 0, student.GroupID)
	require.EqualValues(t, "Two Two", student.Name)
	require.EqualValues(t, "two", student.Alias)
}

func TestStudentFindByName(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find student by name
	student, err := tdb.FindByName("Three Three") // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get student")
	require.EqualValues(t, 3, student.ID)
	require.EqualValues(t, 10, student.GroupID)
	require.EqualValues(t, "Three Three", student.Name)
	require.EqualValues(t, "three", student.Alias)
}

func TestStudentFindByGroup(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find students by group ID
	students, err := tdb.FindByGroup(10) // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get students")
	require.Len(t, students, 2) // this runs after upsert, so there should be 2
}

func TestStudentFindGroups(t *testing.T) {
	// new database type
	tdb := StudentDatabase{}

	// attempt to find students by group ID
	groupIDs, err := tdb.FindGroups()

	// assert output
	require.Nil(t, err, "failed to get group IDs")
	require.Len(t, groupIDs, 1) // student 1 and 3 are in group 10, student 2 doesn't have group (0)
}

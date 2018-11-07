package db

import (
	"testing"

	"hacking-portal/models"

	"github.com/stretchr/testify/require"
)

func TestAnswerUpsert(t *testing.T) {
	// new database type
	tdb := AnswerDatabase{}

	// attempt to insert
	answer, err := tdb.Upsert(models.Answer{1, 1, "foo"})

	// assert output
	require.Nil(t, err, "failed to insert answer")
	require.EqualValues(t, 1, answer.TaskID)
	require.EqualValues(t, 1, answer.GroupID)
	require.EqualValues(t, "foo", answer.Answer)
}

func TestAnswerFindAll(t *testing.T) {
	// new database type
	tdb := AnswerDatabase{}

	// attempt to find answers
	answers, err := tdb.FindAll()

	// assert output
	require.Nil(t, err, "failed to get answers")
	require.Len(t, answers, 1) // this runs after upsert, so there should be 1
}

func TestAnswerFindByID(t *testing.T) {
	// new database type
	tdb := AnswerDatabase{}

	// attempt to find answer by ID
	answer, err := tdb.FindByID(1) // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get single answer")
	require.EqualValues(t, 1, answer.TaskID)
}

func TestAnswerFindByGroup(t *testing.T) {
	// new database type
	tdb := AnswerDatabase{}

	// attempt to find answers by group
	answers, err := tdb.FindByGroup(1) // from the Upsert test

	// assert output
	require.Nil(t, err, "failed to get answers")
	require.Len(t, answers, 1) // this runs after upsert, so there should be 1
}

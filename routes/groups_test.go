package routes

import (
	"bytes"
	"context"
	"hacking-portal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGroups(t *testing.T) {
	sdb := new(mockStudentStorage)
	sdb.Upsert(models.Student{"actual", "Actual User", 1})

	testData := []struct {
		user string
		code int
	}{
		{user: "test", code: http.StatusOK},
		{user: "actual", code: http.StatusTemporaryRedirect},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "session_user_id", data.user))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := GroupsEndpoint{sdb}

		// serve the handler
		handler := http.HandlerFunc(ep.GetGroups)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestPostJoinGroup(t *testing.T) {
	sdb := new(mockStudentStorage)
	sdb.Upsert(models.Student{"ungrouped", "Ungrouped User", 0})
	sdb.Upsert(models.Student{"grouped", "Grouped User", 1})

	testData := []struct {
		body string
		code int
		user string
	}{
		{body: ``, code: http.StatusBadRequest, user: "ungrouped"},
		{body: `{"foo":0}`, code: http.StatusBadRequest, user: "ungrouped"},
		{body: `{"groupID":1}`, code: http.StatusBadRequest, user: "invalid"},
		{body: `{"groupID":1}`, code: http.StatusBadRequest, user: "grouped"},
		{body: `{"groupID":"1"}`, code: http.StatusBadRequest, user: "ungrouped"},
		{body: `{"groupID":0}`, code: http.StatusBadRequest, user: "ungrouped"},
		{body: `{"groupID":1}`, code: http.StatusOK, user: "ungrouped"},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(data.body)))
		req = req.WithContext(context.WithValue(req.Context(), "session_user_id", data.user))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := GroupsEndpoint{sdb}

		// serve the handler
		handler := http.HandlerFunc(ep.PostJoinGroup)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestGroupsRouter(t *testing.T) {
	var r *chi.Mux
	assert.IsType(t, r, GroupsRouter())
}

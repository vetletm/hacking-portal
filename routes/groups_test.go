package routes

import (
	"bytes"
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
		cookie http.Cookie
		code   int
	}{
		{cookie: mockSession("test", false), code: http.StatusOK},
		{cookie: mockSession("actual", true), code: http.StatusTemporaryRedirect},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&data.cookie)

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
		body   string
		code   int
		cookie http.Cookie
	}{
		{body: ``, code: http.StatusBadRequest, cookie: mockSession("ungrouped", true)},
		{body: `{"foo":0}`, code: http.StatusBadRequest, cookie: mockSession("ungrouped", true)},
		{body: `{"groupID":1}`, code: http.StatusBadRequest, cookie: mockSession("invalid", false)},
		{body: `{"groupID":1}`, code: http.StatusBadRequest, cookie: mockSession("grouped", true)},
		{body: `{"groupID":"1"}`, code: http.StatusBadRequest, cookie: mockSession("ungrouped", true)},
		{body: `{"groupID":0}`, code: http.StatusBadRequest, cookie: mockSession("ungrouped", true)},
		{body: `{"groupID":1}`, code: http.StatusOK, cookie: mockSession("ungrouped", true)},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(data.body)))
		req.AddCookie(&data.cookie)

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

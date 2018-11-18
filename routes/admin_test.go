package routes

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminDashboard(t *testing.T) {
	// create a request to pass to the handler
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), "session_user_id", "hei"))

	// create a response recorder to record the response from the handler
	res := httptest.NewRecorder()

	// prepare the endpoint with mocked storage
	ep := AdminEndpoint{
		Machines: new(mockMachineStorage),
		Students: new(mockStudentStorage),
	}

	// serve the handler
	handler := http.HandlerFunc(ep.GetDashboard)
	handler.ServeHTTP(res, req)

	// test the status
	require.Equal(t, http.StatusOK, res.Code, "handler returned wrong status code")
}

func TestPostMachineAssign(t *testing.T) {
	testData := []struct {
		body string
		code int
	}{
		{body: ``, code: http.StatusBadRequest},
		{body: `{"machineUUID":"0000"}`, code: http.StatusBadRequest},
		{body: `{"groupID":0}`, code: http.StatusBadRequest},
		{body: `{"machineUUID":"1111","groupID":1}`, code: http.StatusOK},
		{body: `{"machineUUID":"1111","groupID":-1}`, code: http.StatusInternalServerError},
		{body: `{"machineUUID":"0000","groupID":1}`, code: http.StatusNotFound},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(data.body)))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := AdminEndpoint{
			Machines: new(mockMachineStorage),
			Students: new(mockStudentStorage),
		}

		// serve the handler
		handler := http.HandlerFunc(ep.PostMachineAssign)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestAdminRouter(t *testing.T) {
	var r *chi.Mux
	assert.IsType(t, r, AdminRouter())
}

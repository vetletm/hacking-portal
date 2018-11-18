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

func TestAdminDashboard(t *testing.T) {
	cookie := mockSession("test", true)

	// create a request to pass to the handler
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&cookie)

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
	mdb := new(mockMachineStorage)
	mdb.Upsert(models.Machine{
		Name:    "test1",
		UUID:    "1111",
		GroupID: 1,
		Address: "1.1.1.1",
	})

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
			Machines: mdb,
			Students: new(mockStudentStorage),
		}

		// serve the handler
		handler := http.HandlerFunc(ep.PostMachineAssign)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestAdminMachineRestart(t *testing.T) {
	// create a request to pass to the handler
	req := httptest.NewRequest("POST", "/", nil)

	// create a response recorder to record the response from the handler
	res := httptest.NewRecorder()

	// preprare a new context
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", "0000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// prepare the endpoint with mocked storage
	ep := AdminEndpoint{
		Machines: new(mockMachineStorage),
		Students: new(mockStudentStorage),
	}

	// serve the handler
	handler := http.HandlerFunc(ep.PostMachineRestart)
	handler.ServeHTTP(res, req)

	// test the status
	require.Equal(t, http.StatusNotFound, res.Code, "handler returned wrong status code")
}

func TestAdminRouter(t *testing.T) {
	var r *chi.Mux
	assert.IsType(t, r, AdminRouter())
}

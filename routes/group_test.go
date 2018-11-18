package routes

import (
	"context"
	"hacking-portal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupDashboard(t *testing.T) {
	sdb := new(mockStudentStorage)
	sdb.Upsert(models.Student{"actual", "Actual User", 1})

	mdb := new(mockMachineStorage)
	mdb.Upsert(models.Machine{
		Name:    "test1",
		UUID:    "1111",
		GroupID: 1,
		Address: "1.1.1.1",
	})

	testData := []struct {
		user string
		code int
	}{
		{user: "test", code: http.StatusTemporaryRedirect},
		{user: "actual", code: http.StatusOK},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "session_user_id", data.user))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := GroupEndpoint{
			Machines: mdb,
			Students: sdb,
		}

		// serve the handler
		handler := http.HandlerFunc(ep.GetDashboard)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestGroupMachineRestart(t *testing.T) {
	mdb := new(mockMachineStorage)
	mdb.Upsert(models.Machine{
		Name:    "test1",
		UUID:    "1111",
		GroupID: 2,
		Address: "1.1.1.1",
	})

	sdb := new(mockStudentStorage)
	sdb.Upsert(models.Student{"ungrouped", "Ungrouped User", 0})
	sdb.Upsert(models.Student{"grouped", "Grouped User", 1})

	testData := []struct {
		uuid string
		code int
		user string
	}{
		{uuid: "", code: http.StatusBadRequest, user: "invalid"},
		{uuid: "", code: http.StatusBadRequest, user: "ungrouped"},
		{uuid: "0000", code: http.StatusBadRequest, user: "grouped"},
		{uuid: "1111", code: http.StatusBadRequest, user: "grouped"},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("POST", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "session_user_id", data.user))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := GroupEndpoint{
			Students: sdb,
			Machines: mdb,
		}

		// prepare context
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("uuid", data.uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		// serve the handler
		handler := http.HandlerFunc(ep.PostMachineRestart)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestGetLeaveGroup(t *testing.T) {
	sdb := new(mockStudentStorage)
	sdb.Upsert(models.Student{"ungrouped", "Ungrouped User", 0})
	sdb.Upsert(models.Student{"grouped", "Grouped User", 1})

	testData := []struct {
		code int
		user string
	}{
		{code: http.StatusBadRequest, user: "invalid"},
		{code: http.StatusBadRequest, user: "ungrouped"},
		{code: http.StatusTemporaryRedirect, user: "grouped"},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("POST", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "session_user_id", data.user))

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// prepare the endpoint with mocked storage
		ep := GroupEndpoint{
			Students: sdb,
			Machines: new(mockMachineStorage),
		}

		// serve the handler
		handler := http.HandlerFunc(ep.GetLeaveGroup)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestGroupRouter(t *testing.T) {
	var r *chi.Mux
	assert.IsType(t, r, GroupRouter())
}

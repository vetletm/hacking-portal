package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mockSession(user string, valid bool) http.Cookie {
	if sessions == nil {
		sessions = make(map[string]Session)
	}

	if !valid {
		return http.Cookie{
			Name:    "session_token",
			Value:   user,
			Expires: time.Now(),
		}
	}

	expiration := time.Now().Add(time.Minute)
	cookie := http.Cookie{
		Name:    "session_token",
		Value:   user,
		Expires: expiration,
	}

	sessions[user] = Session{
		Username: user,
		Expires:  expiration,
	}

	return cookie
}

func TestGetLogin(t *testing.T) {
	// create a request to pass to the handler
	req := httptest.NewRequest("GET", "/", nil)

	// create a response recorder to record the response from the handler
	res := httptest.NewRecorder()

	// serve the handler
	handler := http.HandlerFunc(GetLogin)
	handler.ServeHTTP(res, req)

	// test the status
	require.Equal(t, http.StatusOK, res.Code, "handler returned wrong status code")
}

func TestPostLogin(t *testing.T) {
	// create a request to pass to the handler
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("")))

	// create a response recorder to record the response from the handler
	res := httptest.NewRecorder()

	// serve the handler
	handler := http.HandlerFunc(PostLogin)
	handler.ServeHTTP(res, req)

	// test the status
	require.Equal(t, http.StatusBadRequest, res.Code, "handler returned wrong status code")
}

func TestGetLogout(t *testing.T) {
	testData := []struct {
		cookie http.Cookie
		code   int
	}{
		{cookie: http.Cookie{Name: "invalid", Value: "", Expires: time.Now()}, code: http.StatusTemporaryRedirect},
		{cookie: http.Cookie{Name: "session_token", Value: "test", Expires: time.Now().Add(time.Minute)}, code: http.StatusTemporaryRedirect},
	}

	for _, data := range testData {
		// create a request to pass to the handler
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&data.cookie)

		// create a response recorder to record the response from the handler
		res := httptest.NewRecorder()

		// serve the handler
		handler := http.HandlerFunc(GetLogout)
		handler.ServeHTTP(res, req)

		// test the status
		require.Equal(t, data.code, res.Code, "handler returned wrong status code")
	}
}

func TestRedirectLogin(t *testing.T) {
	// create a request to pass to the handler
	req := httptest.NewRequest("GET", "/", nil)

	// create a response recorder to record the response from the handler
	res := httptest.NewRecorder()

	// serve the handler
	handler := http.HandlerFunc(RedirectLogin)
	handler.ServeHTTP(res, req)

	// test the status
	require.Equal(t, http.StatusTemporaryRedirect, res.Code, "handler returned wrong status code")
}

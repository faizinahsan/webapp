package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"personal-projects/webapp/pkg/data"
	"strings"
	"testing"
	"time"
)

func Test_app_authentication(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"not json", `Im Not JSON`, http.StatusUnauthorized},
		{"empty json", `{}`, http.StatusUnauthorized},
		{"empty email", `{"email":""}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@example.com","password":""}`, http.StatusUnauthorized},
		{"invalid password", `{"email":"admin@someotherdomain.com",password":"secret"}`, http.StatusUnauthorized},
	}

	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
func Test_app_refresh(t *testing.T) {
	var tests = []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid token", "", http.StatusOK, true},
		{"valid but not yet ready to expire", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusBadRequest, false},
	}
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}
	oldRefreshTime := refreshTokenExpiry
	for _, e := range tests {
		var tkn string
		if e.token == "" {
			if e.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}
			tokens, _ := app.generateTokenPair(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		postedData := url.Values{
			"refresh_token": {tkn},
		}

		req, _ := http.NewRequest("POST", "/refresh_token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.refresh)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status of %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		refreshTokenExpiry = oldRefreshTime
	}
}

func Test_app_userHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"allUser", "GET", "", "", app.allUsers, http.StatusOK},
		{"deleteUser", "DELETE", "", "1", app.deleteUser, http.StatusNoContent},
		{"deleteUser bad url param", "DELETE", "", "Y", app.deleteUser, http.StatusBadRequest},
		{"getUser valid", "GET", "", "1", app.getUser, http.StatusOK},
		{"getUser invalid", "GET", "", "2", app.getUser, http.StatusBadRequest},
		{"getUser bad url param", "GET", "", "Y", app.getUser, http.StatusBadRequest},
		{
			name:           "updateUser valid",
			method:         "PATCH",
			json:           `{"id":1,"email":"admin@example.com","first_name":"Admin","last_name":"User"}`,
			paramID:        "",
			handler:        app.updateUser,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "updateUser invalid json",
			method:         "PATCH",
			json:           `{"id":1,"email":"admin@example.com",first_name":"Admin","last_name":"User"}`,
			paramID:        "",
			handler:        app.updateUser,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "insertUser valid",
			method:         "PUT",
			json:           `{"email":"admin@example.com","first_name":"Admin","last_name":"User"}`,
			paramID:        "",
			handler:        app.insertUser,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "insertUser invalid",
			method:         "PUT",
			json:           `{"foo":"bar","email":"admin@example.com","first_name":"Admin","last_name":"User"}`,
			paramID:        "",
			handler:        app.insertUser,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "insertUser invalid json",
			method:         "PUT",
			json:           `{email":"admin@example.com","first_name":"Admin","last_name":"User"}`,
			paramID:        "",
			handler:        app.insertUser,
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}
		if e.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatus {
			t.Errorf("%s: expected status %d, but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

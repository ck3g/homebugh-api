package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/ck3g/homebugh-api/pkg/models/mock"
)

func TestAccountsHandler(t *testing.T) {
	app := application{
		models: models.Models{
			AuthSessions: &mock.AuthSessionModel{},
		},
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name           string
		token          string
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:           "With valid token",
			token:          "Bearer valid-token",
			wantStatusCode: http.StatusOK,
			wantBody:       []byte(`{"accounts":[]}`),
		},
		{
			name:           "with blank token",
			token:          "",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"invalid or missing authentication token"}`),
		},
		{
			name:           "with non-bearer token",
			token:          "Notbearer valid-token",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"invalid or missing authentication token"}`),
		},
		{
			name:           "with invalid token",
			token:          "Bearer invalid-token",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"invalid or missing authentication token"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/accounts", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			c := ts.Client()
			rs, err := c.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantStatusCode {
				t.Errorf("want status %d; got %d", tt.wantStatusCode, rs.StatusCode)
			}

			defer rs.Body.Close()
			body, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(body) != string(tt.wantBody) {
				t.Errorf("want body to be equal to `%q`; got `%q`", tt.wantBody, string(body))
			}
		})
	}
}

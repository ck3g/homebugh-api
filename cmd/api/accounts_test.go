package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ck3g/homebugh-api/pkg/jsonh"
	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/ck3g/homebugh-api/pkg/models/mock"
)

func TestAccountsHandler(t *testing.T) {
	app := application{
		models: models.Models{
			Accounts:     &mock.AccountModel{},
			AuthSessions: &mock.AuthSessionModel{},
			Users:        &mock.UserModel{},
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
			wantBody: []byte(`
				{
					"accounts": [
						{
							"id": 1,
							"name": "Bank",
							"balance": 1000,
							"currency": { "id": 1, "name": "Euro", "unit": "€" },
							"status": "active",
							"show_in_summary": true
						},
						{
							"id": 2,
							"name": "Cash",
							"balance": 100.5,
							"currency": { "id": 1, "name": "Euro", "unit": "€" },
							"status": "active",
							"show_in_summary": true
						}
					],
					"metadata":{
						"current_page": 1,
						"page_size": 20,
						"first_page": 1,
						"last_page": 1,
						"total_records":2
					}}`),
		},
		{
			name:           "With valid token of second user",
			token:          "Bearer valid-token-2",
			wantStatusCode: http.StatusOK,
			wantBody: []byte(`
				{
					"accounts": [
						{
							"id": 3,
							"name": "Bank",
							"balance": 500,
							"currency": { "id": 1, "name": "Euro", "unit": "€" },
							"status": "active",
							"show_in_summary": true
						},
						{
							"id": 4,
							"name": "Cash",
							"balance": 30.5,
							"currency": { "id": 1, "name": "Euro", "unit": "€" },
							"status": "active",
							"show_in_summary": true
						}
					],
					"metadata":{
						"current_page": 1,
						"page_size": 20,
						"first_page": 1,
						"last_page": 1,
						"total_records": 2
					}}`),
		},
		{
			name:           "with blank token",
			token:          "",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"you must be authenticated to access this resourse"}`),
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

			if !jsonh.Equal(body, tt.wantBody) {
				t.Errorf("want body to be equal to \n`%q`\ngot \n`%q`\n", tt.wantBody, string(body))
			}
		})
	}
}

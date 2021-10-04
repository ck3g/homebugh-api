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

func TestCategoriesHandler(t *testing.T) {
	app := application{
		models: models.Models{
			Users:        &mock.UserModel{},
			AuthSessions: &mock.AuthSessionModel{},
			Categories:   &mock.CategoryModel{},
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
					"categories": [
						{
							"id": 1,
							"name": "Food",
							"category_type": { "id": 2, "name": "expense" },
							"inactive": false
						}
					],
					"metadata": {
						"current_page": 1,
						"page_size": 20,
						"first_page": 1,
						"last_page": 1,
						"total_records": 1
					}}`),
		},
		{
			name:           "With valid token of second user",
			token:          "Bearer valid-token-2",
			wantStatusCode: http.StatusOK,
			wantBody: []byte(`
				{
					"categories": [
						{
							"id": 2,
							"name":"Groceries",
							"category_type": { "id": 2, "name": "expense" },
							"inactive": false
						}
					],
					"metadata": {
						"current_page": 1,
						"page_size": 20,
						"first_page": 1,
						"last_page": 1,
						"total_records": 1
					}}`),
		},
		{
			name:           "With blank token",
			token:          "",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"you must be authenticated to access this resourse"}`),
		},
		{
			name:           "With non bearer token",
			token:          "Notbearer invalid-token",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"invalid or missing authentication token"}`),
		},
		{
			name:           "With invalid token",
			token:          "Bearer invalid-token",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       []byte(`{"error":"invalid or missing authentication token"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/categories", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			c := ts.Client()
			rs, err := c.Do(req)
			// rs, err := c.Get(ts.URL + "/categories")
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
				t.Errorf("want body to be equal to `%q`; got `%q`", tt.wantBody, string(body))
			}
		})
	}
}

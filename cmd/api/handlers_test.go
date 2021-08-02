package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/ck3g/homebugh-api/pkg/models/mock"
)

func TestHealthHandler(t *testing.T) {
	app := application{}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want status %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	wantBody := `{"status": "OK"}`
	if string(body) != wantBody {
		t.Errorf("want body to be equal to `%q`; got `%q`", wantBody, string(body))
	}
}

func TestCreateTokenHandler(t *testing.T) {
	app := application{
		models: models.Models{
			Users: &mock.UserModel{},
		},
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		body     []byte
		wantCode int
		wantBody []byte
	}{
		{
			name:     "Valid credentials",
			body:     []byte(`{"email": "user@example.com", "password": "password"}`),
			wantCode: http.StatusCreated,
			wantBody: []byte(`{"result": "OK", "token": "valid-token"}`),
		},
		{
			name:     "Valid credentials with spaces",
			body:     []byte(`{"email": " user@example.com ", "password": " password "}`),
			wantCode: http.StatusCreated,
			wantBody: []byte(`{"result": "OK", "token": "valid-token"}`),
		},
		{
			name:     "Not confirmed user",
			body:     []byte(`{"email": "not-confirmed@example.com", "password": "password"}`),
			wantCode: http.StatusUnprocessableEntity,
			wantBody: []byte(`{"result": "Error", "message": "User not confirmed"}`),
		},
		{
			name:     "Empty email",
			body:     []byte(`{"email": "", "password": "password"}`),
			wantCode: http.StatusUnprocessableEntity,
			wantBody: []byte(`{"result": "Error", "message": "User does not exist"}`),
		},
		{
			name:     "Empty password",
			body:     []byte(`{"email": "user@example.com", "password": ""}`),
			wantCode: http.StatusUnprocessableEntity,
			wantBody: []byte(`{"result": "Error", "message": "Invalid credentials"}`),
		},
		{
			name:     "Invalid JSON body",
			body:     []byte(`{"email"`),
			wantCode: http.StatusBadRequest,
			wantBody: []byte("Bad Request\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := ts.Client().Post(ts.URL+"/token", "application/json", bytes.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantCode {
				t.Errorf("want status %d; got %d", tt.wantCode, rs.StatusCode)
			}

			defer rs.Body.Close()
			body, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(body) != string(tt.wantBody) {
				t.Errorf("want body to be equal to `%q`; got `%q`", tt.wantBody, body)
			}
		})
	}

	t.Run("GET request", func(t *testing.T) {
		rs, err := ts.Client().Get(ts.URL + "/token")
		if err != nil {
			t.Fatal(err)
		}

		if rs.StatusCode != http.StatusNotFound {
			t.Errorf("want status %d; got %d", http.StatusNotFound, rs.StatusCode)
		}
	})
}

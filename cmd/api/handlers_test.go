package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHealth(t *testing.T) {
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

func TestCreateToken(t *testing.T) {
	app := application{}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		password string
		wantCode int
		wantBody []byte
	}{
		{"Valid credentials", "user@example.com", "password", http.StatusCreated, []byte(`{"result": "OK"}`)},
		{"Empty email", "", "password", http.StatusUnprocessableEntity, []byte(`{"result": "Error"}`)},
		{"Empty password", "user@example.com", "", http.StatusUnprocessableEntity, []byte(`{"result": "Error"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			rs, err := ts.Client().PostForm(ts.URL+"/token", form)
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

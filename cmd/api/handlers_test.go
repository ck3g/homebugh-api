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

func TestToken(t *testing.T) {
	app := application{}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	form := url.Values{}
	rs, err := ts.Client().PostForm(ts.URL+"/token", form)
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusCreated {
		t.Errorf("want status %d; got %d", http.StatusCreated, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	wantBody := `{"result": "OK"}`
	if string(body) != wantBody {
		t.Errorf("want body to be equal to `%q`; got `%q`", wantBody, string(body))
	}
}

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	app := application{
		metadata: appMetadata{
			version:   "0.0.1+test",
			buildTime: "2021-09-06T14:40:56+0200",
		},
	}

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

	wantBody := `{"buildTime":"2021-09-06T14:40:56+0200","status":"OK","version":"0.0.1+test"}`
	if string(body) != wantBody {
		t.Errorf("want body to be equal to `%q`; got `%q`", wantBody, string(body))
	}
}

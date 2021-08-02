package main

import "net/http"

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "OK",
	}
	// TODO: check for errors
	app.writeJSON(w, http.StatusOK, env, nil)
}

package main

import "net/http"

func (app *application) accountsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Authorization")

	env := envelope{
		"accounts": []string{},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

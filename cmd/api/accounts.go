package main

import (
	"net/http"
	"strings"
)

func (app *application) accountsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Authorization")

	// TODO: extract into middleware
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	token := headerParts[1]
	_, err := app.models.AuthSessions.GetByToken(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	env := envelope{
		"accounts": []string{},
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

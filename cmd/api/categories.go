package main

import (
	"net/http"
	"strings"
)

type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (app *application) categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Authorization")

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
	if token != "valid-token" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	categories := []category{
		{ID: 1, Name: "Food"},
	}
	env := envelope{
		"categories": categories,
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

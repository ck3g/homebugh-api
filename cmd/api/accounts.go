package main

import (
	"net/http"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
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
	session, err := app.models.AuthSessions.GetByToken(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	qs := r.URL.Query()
	filters := models.Filters{
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	accounts, metadata, err := app.models.Accounts.All(session.UserID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"accounts": accounts,
		"metadata": metadata,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

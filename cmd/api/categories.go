package main

import (
	"net/http"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
)

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
	session, err := app.models.AuthSessions.GetByToken(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	qs := r.URL.Query()
	filters := models.Filters{
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "per_page", 20),
	}

	categories, err := app.models.Categories.All(session.UserID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"categories": categories,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

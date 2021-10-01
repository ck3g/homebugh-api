package main

import (
	"net/http"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func (app *application) accountsHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetSession(r)

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

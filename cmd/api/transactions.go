package main

import (
	"net/http"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func (app *application) transactionsHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetSession(r)

	qs := r.URL.Query()
	filters := models.Filters{
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	transactions, metadata, err := app.models.Transactions.All(session.UserID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"transactions": transactions,
		"metadata":     metadata,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

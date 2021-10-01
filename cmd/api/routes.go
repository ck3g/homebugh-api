package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.healthHandler)

	mux.HandleFunc("/token", app.createTokenHandler)

	mux.HandleFunc("/categories", app.requireAuthentication(app.categoriesHandler))

	mux.HandleFunc("/accounts", app.requireAuthentication(app.accountsHandler))

	return app.rateLimit(app.authenticate(mux))
}

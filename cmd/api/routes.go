package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.healthHandler)

	mux.HandleFunc("/token", app.createTokenHandler)

	mux.HandleFunc("/categories", app.categoriesHandler)

	mux.HandleFunc("/accounts", app.accountsHandler)

	return app.rateLimit(mux)
}

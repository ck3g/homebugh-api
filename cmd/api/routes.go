package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/token", app.createTokenHandler)

	return app.rateLimit(mux)
}

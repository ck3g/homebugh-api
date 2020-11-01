package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.health)
	mux.HandleFunc("/token", app.createToken)

	return mux
}

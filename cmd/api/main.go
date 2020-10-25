package main

import (
	"net/http"
)

type application struct{}

func main() {
	app := &application{}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	srv.ListenAndServe()
}

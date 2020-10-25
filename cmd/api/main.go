package main

import (
	"crypto/tls"
	"net/http"
)

type application struct{}

func main() {
	app := &application{}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      ":8080",
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
}

package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "OK"}`))
}

func (app *application) createToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	f := r.PostForm
	email := strings.TrimSpace(f.Get("email"))
	password := strings.TrimSpace(f.Get("password"))

	if email != "user@example.com" || password != "password" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"result": "Error", "message": "Invalid credentials"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	token := "valid-token"
	rBody := fmt.Sprintf(`{"result": "OK", "token": "%s"}`, token)
	w.Write([]byte(rBody))
}

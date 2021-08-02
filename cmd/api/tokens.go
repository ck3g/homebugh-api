package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type createTokenRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.notFoundResponse(w, r)
		return
	}

	var req createTokenRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	token, err := app.models.Users.Authenticate(email, password)
	if err != nil {
		message := createTokenErrorMsg(err)
		app.validationErrorResponse(w, r, map[string]string{"token": message})
		return
	}

	env := envelope{
		"result": "OK",
		"token":  token,
	}

	err = app.writeJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func createTokenErrorMsg(err error) string {
	if errors.Is(err, models.ErrNoRecord) {
		return "User does not exist"
	}

	if errors.Is(err, models.ErrUserNotConfirmed) {
		return "User not confirmed"
	}

	if errors.Is(err, models.ErrWrongPassword) {
		return "Invalid credentials"
	}

	return "Something went wrong"
}

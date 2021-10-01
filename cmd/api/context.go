package main

import (
	"context"
	"net/http"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type contextKey string

const sessionContextKey = contextKey("session")

func (app *application) contextSetSession(r *http.Request, session *models.AuthSession) *http.Request {
	ctx := context.WithValue(r.Context(), sessionContextKey, session)
	return r.WithContext(ctx)
}

func (app *application) contextGetSession(r *http.Request) *models.AuthSession {
	session, ok := r.Context().Value(sessionContextKey).(*models.AuthSession)
	if !ok {
		panic("missing session value in request context")
	}

	return session
}

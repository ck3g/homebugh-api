package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.metadata.environment == "production" {
			ip := realip.FromRequest(r)

			mu.Lock()

			// TODO: extract into a config
			rps := 2
			burst := 4

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(rps), burst),
				}
			}

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			// app.invalidAuthenticationTokenResponse(w, r)
			r = app.contextSetSession(r, models.AnonymousSession)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		session, err := app.models.AuthSessions.GetByToken(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		r = app.contextSetSession(r, session)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := app.contextGetSession(r)

		if session.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package main

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

// standardMiddleware defines middleware chain for all routes.
func (app *application) standardMiddleware() alice.Chain {
	return alice.New(app.recoverPanic, app.logRequest, commonHeaders)
}

// dynamicMiddleware defines middleware chain for dynamic routes.
func (app *application) dynamicMiddleware() alice.Chain {
	return alice.New(app.noSurf)
}

// commonHeaders sets headers for all requests.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Add("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "deny")
		w.Header().Add("X-XSS-Protection", "0")

		w.Header().Add("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a catch-all for panics that occur within the go routine the request is being handled in.
// This allows us to present a 500 response to the client in the event of an un-recovered panic during request handling.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.logger.Error(
					"recovered from panic",
					slog.String("ip", r.RemoteAddr),
					slog.String("method", r.Method),
					slog.String("uri", r.URL.RequestURI()),
					slog.Any("error", err),
				)
				w.Write([]byte("500 - Internal Server Error"))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// logRequest logs anything request related to the application logger.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info(
			"received request",
			slog.String("ip", ip),
			slog.String("proto", proto),
			slog.String("method", method),
			slog.String("uri", uri),
		)

		next.ServeHTTP(w, r)
	})
}

// noSurf configures the noSurf CSRFHandler to manage csrf tokens.
func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",

		// dependent on if we are running with tls or not.
		Secure: app.config.tls.enabled,
	})

	return csrfHandler
}

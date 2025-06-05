package main

import (
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"full-stack-demo.tharris.uk/ui"
)

// routes returns a http.Handler with the application routes registered.
func (app *application) routes() (http.Handler, error) {
	router := httprouter.New()

	// register the http problem handlers
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("405 - Method Not Allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// register static routes
	dir, err := fs.Sub(ui.Files, "static")
	if err != nil {
		return nil, err
	}
	router.ServeFiles("/static/*filepath", http.FS(dir))

	router.HandlerFunc(http.MethodGet, "/ping", app.pingHandler)

	// register the dynamic routes with the dynamic middleware chain
	dynamicMiddleware := app.dynamicMiddleware()
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.homeHandler))

	// all routes are wrapped in the standard middleware chain
	return app.standardMiddleware().Then(router), nil
}

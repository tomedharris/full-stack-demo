package main

import (
	"net/http"
)

// homeHandler is the HTTP GET handler for the application home page.
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

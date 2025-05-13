package routes

import (
	"net/http"
)

func RegisterFileRoutes(mux *http.ServeMux) {
	// serve static files
	fileserver := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
}

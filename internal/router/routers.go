package router

import (
	"net/http"

	"github.com/smrrazavian/url-shortener/internal/router/handlers"
)

// TODO: Complete doc
func AddRoutes(
	mux *http.ServeMux,
) {
	mux.HandleFunc("/save", handlers.SaveURLHandler)
	mux.HandleFunc("/", handlers.GetURLHandler)
}

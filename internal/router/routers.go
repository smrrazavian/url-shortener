package router

import (
	"net/http"

	"github.com/smrrazavian/url-shortener/internal/router/handlers"
	"github.com/smrrazavian/url-shortener/internal/router/middleware"
)

// TODO: Complete doc
func AddRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("/save", middleware.JwtMiddleware(http.HandlerFunc(handlers.SaveURLHandler)))
	mux.HandleFunc("/", handlers.GetURLHandler)
}

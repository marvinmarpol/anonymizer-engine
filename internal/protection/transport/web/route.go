package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Route interface {
	Deidentify(w http.ResponseWriter, r *http.Request)
	Reidentify(w http.ResponseWriter, r *http.Request)
	GetCypher(w http.ResponseWriter, r *http.Request)
	RotateKeys(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(si Route, Middlewares []MiddlewareFunc) http.Handler {
	r := chi.NewRouter()
	for _, middleware := range Middlewares {
		r.Use(middleware)
	}

	r.Post("/deidentify", si.Deidentify)
	r.Post("/reidentify", si.Reidentify)
	r.Post("/rotate", si.RotateKeys)

	return r
}

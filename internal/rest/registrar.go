package rest

import "github.com/go-chi/chi/v5"

type RouteRegistrar interface {
	RegisterRoutes(router *chi.Mux)
}

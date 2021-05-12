package main

import "github.com/go-chi/chi/v5"

func InitRoutes() {
	// Home route
	srv.Router.Get("/", Home)

	// Customers
	srv.Router.Route("/users", func(r chi.Router) {
		r.Get("/", GetUsers)
	})
}

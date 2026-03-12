package api

import (
	"github.com/UGRORF/price-tracker/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func RouterInit() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.MainHandler)
	r.Get("/user/{id}", handlers.GetUser)
	r.Get("/users", handlers.GetAllUsers)
	r.Get("/store/{id}", handlers.GetStore)
	r.Get("/stores", handlers.GetAllStores)
	r.Get("/product/{id}", handlers.GetProduct)
	r.Get("/products", handlers.GetAllProducts)
	r.Get("/offer/{id}", handlers.GetOffer)

	return r
}

package api

import (
	"github.com/UGRORF/price-tracker/internal/api/handlers"
	"github.com/UGRORF/price-tracker/internal/api/middleware"
	"github.com/UGRORF/price-tracker/internal/service/auth"
	"github.com/go-chi/chi/v5"
)

func RouterInit(authHandler *handlers.AuthHandler,
	jwtService *auth.JWTService) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))

		r.Get("/api/stores", handlers.GetAllStores)
		r.Get("/api/store/{id}", handlers.GetStore)
		r.Get("/api/products", handlers.GetAllProducts)
		r.Get("/api/product/{id}", handlers.GetProduct)
		r.Get("/api/offers", handlers.GetOffers)
		r.Get("/api/offer/{id}", handlers.GetOffer)
	})

	return r
}

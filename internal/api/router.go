package api

import (
	"github.com/UGRORF/price-tracker/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func RouterInit() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.MainHandler)

	return r
}

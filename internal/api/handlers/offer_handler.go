package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var offerRepo *postgres.OfferRepo

func InitOfferRepo(repo *postgres.OfferRepo) {
	offerRepo = repo
}

func GetOffer(w http.ResponseWriter, r *http.Request) {
	if offerRepo == nil {
		http.Error(w, "Store repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	offer, err := offerRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Offer not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(offer)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

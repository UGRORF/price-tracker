package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/api/middleware"
	"github.com/UGRORF/price-tracker/internal/domain"
	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var storeRepo *postgres.StoreRepo

func InitStoreRepo(repo *postgres.StoreRepo) {
	storeRepo = repo
}

func GetStore(w http.ResponseWriter, r *http.Request) {
	if storeRepo == nil {
		http.Error(w, "Store repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	store, err := storeRepo.GetById(id)
	if err != nil {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(store)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetAllStores(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if role != string(domain.RoleAdmin) {
		http.Error(w, `{"error": "forbidden: admin rights required"}`, http.StatusForbidden)
		return
	}

	if storeRepo == nil {
		http.Error(w, "Store repository not initialized", http.StatusInternalServerError)
		return
	}

	stores, err := storeRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get stores: "+err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(stores)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

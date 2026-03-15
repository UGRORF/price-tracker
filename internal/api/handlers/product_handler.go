package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var productRepo *postgres.ProductRepo

func InitProductRepo(repo *postgres.ProductRepo) {
	productRepo = repo
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if productRepo == nil {
		http.Error(w, "Product repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	product, err := productRepo.GetById(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if productRepo == nil {
		http.Error(w, "Product repository not initialized", http.StatusInternalServerError)
		return
	}

	products, err := productRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(products)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

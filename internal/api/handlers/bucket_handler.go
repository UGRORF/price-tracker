package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/api/middleware"
	"github.com/UGRORF/price-tracker/internal/domain"
	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var bucketRepo *postgres.BucketRepo

func InitBucketRepo(repo *postgres.BucketRepo) {
	bucketRepo = repo
}

func GetBucketByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if bucketRepo == nil {
		http.Error(w, "Bucket repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	bucket, err := bucketRepo.GetByUserID(id, userID)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(bucket)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetBucket(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if role != string(domain.RoleAdmin) {
		http.Error(w, `{"error": "forbidden: admin rights required"}`, http.StatusForbidden)
		return
	}

	if bucketRepo == nil {
		http.Error(w, "Bucket repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	bucket, err := bucketRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(bucket)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetAllBucketsByUserID(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if bucketRepo == nil {
		http.Error(w, "Bucket repository not initialized", http.StatusInternalServerError)
		return
	}

	buckets, err := bucketRepo.GetAllByUserID(id)
	if err != nil {
		http.Error(w, "Buckets not founds", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(buckets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetAllBuckets(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if role != string(domain.RoleAdmin) {
		http.Error(w, `{"error": "forbidden: admin rights required"}`, http.StatusForbidden)
		return
	}

	if bucketRepo == nil {
		http.Error(w, "Bucket repository not initialized", http.StatusInternalServerError)
		return
	}

	buckets, err := bucketRepo.GetAll()
	if err != nil {
		http.Error(w, "Buckets not founds", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(buckets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func AddBucket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if bucketRepo == nil || productRepo == nil || storeRepo == nil || offerRepo == nil {
		http.Error(w, "repositories is not initialized", http.StatusInternalServerError)
		return
	}

	var req struct {
		ProductName        string  `json:"productName"`
		ProductDescription string  `json:"productDescription,omitempty"`
		StoreURL           string  `json:"storeURL"`
		StoreName          string  `json:"storeName"`
		TargetPrice        float64 `json:"targetPrice,omitempty"`
	}

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProductName == "" {
		http.Error(w, `{"error": "productName is required"}`, http.StatusBadRequest)
		return
	}
	if req.StoreURL == "" {
		http.Error(w, `{"error": "storeURL is required"}`, http.StatusBadRequest)
		return
	}
	if req.StoreName == "" {
		http.Error(w, `{"error": "storeName is required"}`, http.StatusBadRequest)
		return
	}

	product, err := productRepo.GetByName(req.ProductName)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if product == nil {
		product = &domain.Product{}
		product.Name = req.ProductName
		product.Description = req.ProductDescription
		err = productRepo.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	store, err := storeRepo.GetByURL(req.StoreURL)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if store == nil {
		store = &domain.Store{
			Name: req.StoreName,
			URL:  req.StoreURL,
		}
		err = storeRepo.Create(store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	offer := &domain.Offer{
		ProductID: product.ID,
		StoreID:   store.ID,
		//TODO: переделать на минимальную цену по продукту, когда сделаю парсинг цен
		Price:   req.TargetPrice,
		Product: product,
		Store:   store,
	}
	err = offerRepo.Create(offer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucket := &domain.Bucket{
		UserID:      userID,
		OfferID:     offer.ID,
		TargetPrice: req.TargetPrice,
		Offer:       offer,
	}
	err = bucketRepo.Create(bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := json.Marshal(bucket)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

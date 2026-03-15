package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var bucketRepo *postgres.BucketRepo

func InitBucketRepo(repo *postgres.BucketRepo) {
	bucketRepo = repo
}

func GetBucket(w http.ResponseWriter, r *http.Request) {
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

func GetBuckets(w http.ResponseWriter, r *http.Request) {
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
	if bucketRepo == nil {
		http.Error(w, "Bucket repository not initialized", http.StatusInternalServerError)
		return
	}

}

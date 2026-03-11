package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

var userRepo *postgres.UserRepo

func InitUserRepo(repo *postgres.UserRepo) {
	userRepo = repo
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if userRepo == nil {
		http.Error(w, "User repository not initialized", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID must be is int", http.StatusBadRequest)
		return
	}

	user, err := userRepo.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if userRepo == nil {
		http.Error(w, "User repository not initialized", http.StatusInternalServerError)
		return
	}

	users, err := userRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with serialization in JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

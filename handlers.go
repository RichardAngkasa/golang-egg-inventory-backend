package main

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
	db   *Database
	auth *Auth
}

func NewHandlers(db *Database, auth *Auth) *Handlers {
	return &Handlers{db: db, auth: auth}
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.FindUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.auth.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handlers) CreateBulkEggRacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := h.auth.ValidateJWT(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var counts struct {
		BigCount    int `json:"big_count"`
		MediumCount int `json:"medium_count"`
		SmallCount  int `json:"small_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&counts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.db.CreateBulkEggRacks(username, counts.BigCount, counts.MediumCount, counts.SmallCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handlers) GetAllEggRacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, err := h.auth.ValidateJWT(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	racks := h.db.GetAllEggRacks(username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(racks)
}

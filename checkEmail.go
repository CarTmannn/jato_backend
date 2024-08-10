package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type CheckRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Response struct {
	EmailExists    bool `json:"email"`
	UsernameExists bool `json:"username"`
}

func checkEmailAndUsernameExists(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	var response Response

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if email exists in user or builder table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	response.EmailExists = count > 0

	err = db.QueryRow("SELECT COUNT(*) FROM builder WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	response.EmailExists = response.EmailExists || count > 0

	// Check if username exists in user or builder table
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	response.UsernameExists = count > 0

	err = db.QueryRow("SELECT COUNT(*) FROM builder WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	response.UsernameExists = response.UsernameExists || count > 0

	// Return response
	json.NewEncoder(w).Encode(response)
}

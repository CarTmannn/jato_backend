package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type EmailRequest struct {
	Email string `json:"email"`
}

type Response struct {
	Exists bool `json:"exists"`
}

func checkEmailExists(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest
	var response Response

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if email exists in user table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		response.Exists = true
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if email exists in builder table
	err = db.QueryRow("SELECT COUNT(*) FROM builder WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		response.Exists = true
	} else {
		response.Exists = false
	}

	// Return response
	json.NewEncoder(w).Encode(response)
}

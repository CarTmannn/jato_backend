package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type RequestBody struct {
	Email string `json:"email"`
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deleteUserQuery := "DELETE FROM user WHERE email = ?"
	resultUser, err := db.Exec(deleteUserQuery, reqBody.Email)
	if err != nil {
		http.Error(w, "Failed to delete from users table", http.StatusInternalServerError)
		return
	}

	deleteBuilderQuery := "DELETE FROM builder WHERE email = ?"
	resultBuilder, err := db.Exec(deleteBuilderQuery, reqBody.Email)
	if err != nil {
		http.Error(w, "Failed to delete from builders table", http.StatusInternalServerError)
		return
	}

	rowsAffectedUser, err := resultUser.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking affected rows in users table", http.StatusInternalServerError)
		return
	}

	rowsAffectedBuilder, err := resultBuilder.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking affected rows in builders table", http.StatusInternalServerError)
		return
	}

	if rowsAffectedUser == 0 && rowsAffectedBuilder == 0 {
		http.Error(w, "No user found with that email in both tables", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User with email %s deleted from users and builders tables", reqBody.Email)
	}
}

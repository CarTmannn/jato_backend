package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Struktur data untuk User dan Builder
type UserData struct {
	ID                   int    `json:"id"`
	FullName             string `json:"nama_lengkap"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	PlaceOfBirth         string `json:"tempat_lahir"`
	DateOfBirth          string `json:"tanggal_lahir"`
	Gender               string `json:"jenis_kelamin"`
	PhoneNumber          string `json:"no_handphone"`
	NationalID           string `json:"nomor_induk_kependudukan"`
	Occupation           string `json:"pekerjaan"`
	Address              string `json:"alamat"`
	RT                   string `json:"rt"`
	RW                   string `json:"rw"`
	Province             string `json:"provinsi"`
	CityOrRegency        string `json:"kota_kabupaten"`
	District             string `json:"kecamatan"`
	VillageOrSubdistrict string `json:"kelurahan_desa"`
	Role                 string `json:"role"`
	Password             string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var credentials UserData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	user, err := validateCredentials(credentials.Email, credentials.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid credentials")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error validating credentials")
			log.Println(err)
		}
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshalling user data")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func validateCredentials(email string, password string) (*UserData, error) {
	var user UserData

	row := db.QueryRow(`
		SELECT id, nama_lengkap, username, email, tempat_lahir, tanggal_lahir, 
		       jenis_kelamin, no_handphone, nomor_induk_kependudukan, pekerjaan, 
		       alamat, rt, rw, provinsi, kota_kabupaten, kecamatan, kelurahan_desa, 
		       role, password
		FROM user 
		WHERE email = ? AND password = ?`, email, password)

	err := row.Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.PlaceOfBirth,
		&user.DateOfBirth, &user.Gender, &user.PhoneNumber, &user.NationalID, &user.Occupation,
		&user.Address, &user.RT, &user.RW, &user.Province, &user.CityOrRegency, &user.District,
		&user.VillageOrSubdistrict, &user.Role, &user.Password)

	if err == nil {
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	row = db.QueryRow(`
		SELECT id, nama_lengkap, username, email, tempat_lahir, tanggal_lahir, 
		       jenis_kelamin, no_handphone, nomor_induk_kependudukan, pekerjaan, 
		       alamat, rt, rw, provinsi, kota_kabupaten, kecamatan, kelurahan_desa, 
		       role, password
		FROM builder 
		WHERE email = ? AND password = ?`, email, password)

	err = row.Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.PlaceOfBirth,
		&user.DateOfBirth, &user.Gender, &user.PhoneNumber, &user.NationalID, &user.Occupation,
		&user.Address, &user.RT, &user.RW, &user.Province, &user.CityOrRegency, &user.District,
		&user.VillageOrSubdistrict, &user.Role, &user.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

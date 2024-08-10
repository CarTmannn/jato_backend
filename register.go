package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	NamaLengkap            string `json:"nama_lengkap"`
	Username               string `json:"username"`
	Email                  string `json:"email"`
	Password               string `json:"password"`
	TempatLahir            string `json:"tempat_lahir"`
	TanggalLahir           string `json:"tanggal_lahir"`
	JenisKelamin           string `json:"jenis_kelamin"`
	NoHandphone            string `json:"no_handphone"`
	NomorIndukKependudukan string `json:"nomor_induk_kependudukan"`
	Pekerjaan              string `json:"pekerjaan"`
	Alamat                 string `json:"alamat"`
	RT                     string `json:"rt"`
	RW                     string `json:"rw"`
	Provinsi               string `json:"provinsi"`
	KotaKabupaten          string `json:"kota_kabupaten"`
	Kecamatan              string `json:"kecamatan"`
	KelurahanDesa          string `json:"kelurahan_desa"`
	Role                   string `json:"role"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	var err error

	switch user.Role {
	case "builder":
		_, err = db.Exec(`
			INSERT INTO builder (
				nama_lengkap, username, email, password, tempat_lahir, tanggal_lahir, 
				jenis_kelamin, no_handphone, nomor_induk_kependudukan, pekerjaan, 
				alamat, rt, rw, provinsi, kota_kabupaten, kecamatan, kelurahan_desa, role
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.NamaLengkap, user.Username, user.Email, user.Password, user.TempatLahir, user.TanggalLahir,
			user.JenisKelamin, user.NoHandphone, user.NomorIndukKependudukan, user.Pekerjaan,
			user.Alamat, user.RT, user.RW, user.Provinsi, user.KotaKabupaten, user.Kecamatan,
			user.KelurahanDesa, user.Role)
	case "user":
		_, err = db.Exec(`
			INSERT INTO user (
				nama_lengkap, username, email, password, tempat_lahir, tanggal_lahir, 
				jenis_kelamin, no_handphone, nomor_induk_kependudukan, pekerjaan, 
				alamat, rt, rw, provinsi, kota_kabupaten, kecamatan, kelurahan_desa, role
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.NamaLengkap, user.Username, user.Email, user.Password, user.TempatLahir, user.TanggalLahir,
			user.JenisKelamin, user.NoHandphone, user.NomorIndukKependudukan, user.Pekerjaan,
			user.Alamat, user.RT, user.RW, user.Provinsi, user.KotaKabupaten, user.Kecamatan,
			user.KelurahanDesa, user.Role)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid role specified")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting data into database: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User added successfully")
}



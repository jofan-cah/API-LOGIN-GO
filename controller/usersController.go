package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jofan-cah/login-api/database"
	"github.com/jofan-cah/login-api/models"
	"github.com/jofan-cah/login-api/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Secret key untuk JWT

// Register user handler
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	db := database.DB
	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
	_, err = db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login user handler
func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	db := database.DB
	query := "SELECT id, username, password FROM users WHERE username = ?"
	row := db.QueryRow(query, creds.Username)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Get all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	query := "SELECT id, username, email, created_at, updated_at FROM users"

	// Ambil data pengguna dari database
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Error fetching users from database", http.StatusInternalServerError)
		log.Println("Error fetching users:", err)
		return
	}
	defer rows.Close()

	// Iterasi hasil query dan pindai hasilnya
	for rows.Next() {
		var user models.User
		// Pindai semua 5 kolom sesuai dengan urutan yang ada di query
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			log.Println("Error scanning user data:", err)
			return
		}
		users = append(users, user)
	}

	// Pastikan tidak ada error saat iterasi rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		log.Println("Error reading rows:", err)
		return
	}

	// Kirim response dengan data pengguna dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Delete user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	db := database.DB
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

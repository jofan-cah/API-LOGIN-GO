package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/jofan-cah/login-api/controller"
	"github.com/jofan-cah/login-api/models"
	"github.com/jofan-cah/login-api/utils"
)

func TestRegister(t *testing.T) {
	// Persiapan router
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register).Methods("POST")

	// Data pengguna untuk testing
	testCases := []struct {
		name           string
		userData       map[string]string
		expectedStatus int
	}{
		{
			name: "Successful Registration",
			userData: map[string]string{
				"username": "testuser",
				"password": "testpassword",
				"email":    "testuser@example.com",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Input",
			userData: map[string]string{
				"username": "",
				"password": "",
				"email":    "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Konversi data pengguna ke JSON
			jsonData, _ := json.Marshal(tc.userData)

			// Buat request
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Rekam respons
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Verifikasi status code
			assert.Equal(t, tc.expectedStatus, rr.Code,
				"Expected status code %d, got %d", tc.expectedStatus, rr.Code)

			// Jika registrasi berhasil, periksa struktur respons
			if tc.expectedStatus == http.StatusCreated {
				var user models.User
				err := json.Unmarshal(rr.Body.Bytes(), &user)
				assert.NoError(t, err, "Should be able to unmarshal response")
				assert.NotEmpty(t, user.Username, "Username should not be empty")
				assert.NotEmpty(t, user.Email, "Email should not be empty")
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Persiapan router
	r := mux.NewRouter()
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	testCases := []struct {
		name           string
		loginData      map[string]string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Successful Login",
			loginData: map[string]string{
				"username": "testuser",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Invalid Credentials",
			loginData: map[string]string{
				"username": "wronguser",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Konversi data login ke JSON
			jsonData, _ := json.Marshal(tc.loginData)

			// Buat request login
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Rekam respons
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Verifikasi status code
			assert.Equal(t, tc.expectedStatus, rr.Code,
				"Expected status code %d, got %d", tc.expectedStatus, rr.Code)

			// Jika login berhasil, periksa token
			if tc.expectToken {
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err, "Should be able to unmarshal response")

				token := response["token"]
				assert.NotEmpty(t, token, "Token should not be empty")

				// Validasi token (opsional)
				claims, err := utils.ValidateJWT(token)
				assert.NoError(t, err, "Token should be valid")
				assert.NotNil(t, claims, "Claims should not be nil")
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	// Persiapan router
	r := mux.NewRouter()
	r.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")

	// Buat request GET
	req, _ := http.NewRequest("GET", "/users", nil)

	// Rekam respons
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Verifikasi status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	// Parse respons JSON
	var users []models.User
	err := json.Unmarshal(rr.Body.Bytes(), &users)
	assert.NoError(t, err, "Should be able to unmarshal users")

	// Verifikasi setidaknya ada satu pengguna
	assert.True(t, len(users) > 0, "Should have at least one user")

	// Verifikasi struktur pengguna
	for _, user := range users {
		assert.NotEmpty(t, user.ID, "User ID should not be empty")
		assert.NotEmpty(t, user.Username, "Username should not be empty")
		assert.NotEmpty(t, user.Email, "Email should not be empty")
	}
}

func TestDeleteUser(t *testing.T) {
	// Persiapan router
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

	testCases := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Successful Deletion",
			userID:         "1", // Ganti dengan ID pengguna yang valid di database Anda
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User Not Found",
			userID:         "9999", // ID yang tidak ada
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Buat request DELETE
			req, _ := http.NewRequest("DELETE", "/users/"+tc.userID, nil)

			// Rekam respons
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Verifikasi status code
			assert.Equal(t, tc.expectedStatus, rr.Code,
				"Expected status code %d, got %d", tc.expectedStatus, rr.Code)

			// Jika penghapusan berhasil, periksa respons
			if tc.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err, "Should be able to unmarshal response")
				assert.Equal(t, "User deleted successfully", response["message"],
					"Should have success message")
			}
		})
	}
}

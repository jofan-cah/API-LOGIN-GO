package utils

import (
	"errors"
	"os"
	"time"

	"github.com/jofan-cah/login-api/models"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateJWT generates a JWT token for the user
func GenerateJWT(user models.User) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT - Fungsi untuk memvalidasi token JWT
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse token dengan secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma yang digunakan adalah HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtKey, nil
	})

	// Jika parsing gagal atau token tidak valid
	if err != nil {
		return nil, err
	}

	// Casting claims ke tipe yang sesuai
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

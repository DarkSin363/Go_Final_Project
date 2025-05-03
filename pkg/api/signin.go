package api

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("Yandex_Practicum_Final_Project") 

func signinHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJsonError(w, "Authentication not configured", http.StatusInternalServerError)
		return
	}

	if req.Password != expectedPassword {
		writeJsonError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(8 * time.Hour).Unix(),
		"pwd_hash": hashPassword(expectedPassword),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	writeJson(w, map[string]string{"token": tokenString})
}

func hashPassword(pwd string) string {

	return fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
}
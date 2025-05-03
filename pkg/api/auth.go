package api

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        if os.Getenv("TODO_PASSWORD") == "" {
            next.ServeHTTP(w, r)
            return
        }

        if r.URL.Path == "/api/signin" {
            next.ServeHTTP(w, r)
            return
        }

        cookie, err := r.Cookie("token")
		log.Printf("%v",cookie)
        if err != nil {
            writeJsonError(w, "Authentication required", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            writeJsonError(w, "Invalid token", http.StatusUnauthorized)
            return
        }

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			currentPwdHash := hashPassword(os.Getenv("TODO_PASSWORD"))
			if claims["pwd_hash"] != currentPwdHash {
				writeJsonError(w, "Password changed", http.StatusUnauthorized)
				return
			}
		}

        next.ServeHTTP(w, r)
    })
}
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var authConfig map[string]map[string]map[string]bool

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordValidity(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			log.Println("Auth token is missing")
			w.WriteHeader(http.StatusUnauthorized)
			response := Response{
				Ok:      false,
				Message: "Auth token is missing",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		jwtToken := authHeader[1]
		decodedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})

		claims := decodedToken.Claims.(jwt.MapClaims)
		ctx = context.WithValue(ctx, "role", claims["role"].(string))

		if err != nil || !decodedToken.Valid {
			response := Response{
				Ok:      false,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)

		route, _ := mux.CurrentRoute(r).GetPathTemplate()
		hasPermission := authConfig[role][route][r.Method]

		if !hasPermission {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GenerateToken(user *User) (string, error) {
	jwt_expiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Duration(jwt_expiration) * time.Millisecond).Unix(),
		"role": user.Role,
	})

	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

func InitAuthorization() {
	file, _ := os.Open("./authorization.json")

	json.NewDecoder(file).Decode(&authConfig)
}

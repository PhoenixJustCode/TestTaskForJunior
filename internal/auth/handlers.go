package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Регистрация
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		hash, _ := HashPassword(creds.Password)
		_, err = db.Exec("INSERT INTO users(username, password_hash) VALUES($1, $2)", creds.Username, hash)
		if err != nil {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user registered"))
	}
}

// Логин
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		var userID int
		var hash string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", creds.Username).Scan(&userID, &hash)
		if err != nil || !CheckPasswordHash(creds.Password, hash) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		access, refresh, _ := GenerateTokens(userID)
		SaveRefreshToken(db, userID, refresh, time.Now().Add(7*24*time.Hour))

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	}
}

// Обновление токена
func RefreshHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Refresh string `json:"refresh_token"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		userID, err := CheckRefreshToken(db, body.Refresh)
		if err != nil {
			http.Error(w, "invalid refresh", http.StatusUnauthorized)
			return
		}

		access, refresh, _ := GenerateTokens(userID)
		SaveRefreshToken(db, userID, refresh, time.Now().Add(7*24*time.Hour))

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	}
}

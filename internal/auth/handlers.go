package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// структура для регистрации
type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// структура для логина
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 📌 Регистрация
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		// хэшируем пароль
		hash, err := HashPassword(input.Password)
		if err != nil {
			http.Error(w, "hashing error", http.StatusInternalServerError)
			return
		}

		// проверяем уникальность email/username
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)", input.Email, input.Username).Scan(&exists)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

		// сохраняем юзера
		var userID int
		err = db.QueryRow(
			"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
			input.Username, input.Email, hash,
		).Scan(&userID)
		if err != nil {
			http.Error(w, "db insert error", http.StatusInternalServerError)
			return
		}

		// выдаём токены
		getTokens(w,r,db, userID)
	}
}

// 📌 Логин
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		var userID int
		var hash string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", input.Username).Scan(&userID, &hash)
		if err != nil || !CheckPasswordHash(input.Password, hash) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		getTokens(w,r,db, userID)
	}
}

// 📌 Обновление токена
func RefreshHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshToken string

		var body struct {
			Refresh string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.Refresh != "" {
			refreshToken = body.Refresh
		}

		if refreshToken == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if refreshToken == "" {
			http.Error(w, "missing refresh token", http.StatusBadRequest)
			return
		}

		// Проверка refresh токена
		userID, err := CheckRefreshToken(db, refreshToken)
		if err != nil {
			http.Error(w, "invalid refresh", http.StatusUnauthorized)
			return
		}

		getTokens(w,r,db, userID)
	}
}


// функция выхода и удаление 
func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshToken string
		if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
			refreshToken = cookie.Value
		}

		if refreshToken == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if refreshToken != "" {
			if _, err := db.Exec("DELETE FROM refresh_tokens WHERE token=$1", refreshToken); err != nil {
				http.Error(w, "failed to delete token", http.StatusInternalServerError)
				return
			}
		}

		clearCookieHandler(w, r)

		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(`{"message":"logged out"}`))
	}
}



// Генерация новых токенов
func getTokens(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	access, refresh, _ := GenerateTokens(userID)
	
	SaveRefreshToken(db, userID, refresh, time.Now().Add(7*24*time.Hour))

	setCookieHandler(w, r, refresh)
	
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  access,
	})

}




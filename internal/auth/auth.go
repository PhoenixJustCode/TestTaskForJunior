package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func InitJWT(secret string) {
	jwtKey = []byte(secret)
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Хэширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Генерация access + refresh токенов
func GenerateTokens(userID int) (string, string, error) {
	// access токен
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// refresh токен
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExp.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

// Валидация токена
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// Сохраняем refresh в БД
func SaveRefreshToken(db *sql.DB, userID int, token string, exp time.Time) error {
	_, err := db.Exec(
		"INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES($1, $2, $3)",
		userID, token, exp,
	)
	return err
}

// Проверяем refresh в БД
func CheckRefreshToken(db *sql.DB, token string) (int, error) {
	var userID int
	var expiresAt time.Time
	err := db.QueryRow(
		"SELECT user_id, expires_at FROM refresh_tokens WHERE token=$1", token,
	).Scan(&userID, &expiresAt)

	if err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("refresh expired")
	}
	return userID, nil
}

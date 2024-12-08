package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/petersizovdev/MEDODS-T.git/pkg/env"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func generateRefreshToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "generateRefreshToken error:", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func generateAccessToken(userID string, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": userID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	err := env.LoadEnv(".env")
	if err != nil {
		fmt.Println("Err to load .env", err)
	}
	secretString := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secretString))
	if err != nil {
		return "generateAccessToken error:", err
	}
	return tokenString, nil
}

func (a *AuthHandler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		fmt.Println("No user found")
		return
	}

	ip := r.RemoteAddr

	accessToken, err := generateAccessToken(userID, ip)
	if err != nil {
		fmt.Println("Access token generation error:", err)
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		fmt.Println("Refresh token generation error:", err)
		return
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Refresh token hashing error:", err)
		return
	}

	_, err = a.DB.Exec("INSERT INTO refresh_tokens (user_id, token_hash) VALUES ($1, $2)", userID, hashedRefreshToken)
	if err != nil {
		fmt.Println("Failed to save refresh_token:", err)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	refreshToken := r.URL.Query().Get("refresh_token")
	if accessToken == "" || refreshToken == "" {
		fmt.Println("No access_token or refresh_token found")
		return
	}

	err := env.LoadEnv(".env")
	if err != nil {
		fmt.Println("Err to load .env", err)
	}
	secretString := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing error")
		}
		return []byte(secretString), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Invalid access token:", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Invalid access token claims")
		return
	}

	userID := claims["user_id"].(string)
	ip := claims["ip"].(string)
	currentIP := r.RemoteAddr
	if ip != currentIP {
		fmt.Println("IP Адресс изменился, уведомление отправлено на почту пользователя:", userID)
	}

	var storedTokenHash []byte
	err = a.DB.QueryRow("SELECT token_hash FROM refresh_tokens WHERE user_id = $1", userID).Scan(&storedTokenHash)
	if err != nil {
		fmt.Println("Invalid refresh token:", err)
		return
	}

	err = bcrypt.CompareHashAndPassword(storedTokenHash, []byte(refreshToken))
	if err != nil {
		fmt.Println("Invalid refresh token:", err)
		return
	}

	newAccessToken, err := generateAccessToken(userID, currentIP)
	if err != nil {
		fmt.Println("Failed to generate new access token:", err)
		return
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		fmt.Println("Failed to generate new refresh token:", err)
		return
	}

	newHashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash new refresh token:", err)
		return
	}

	_, err = a.DB.Exec("UPDATE refresh_tokens SET token_hash = $1 WHERE user_id = $2", newHashedRefreshToken, userID)
	if err != nil {
		fmt.Println("Failed to update refresh token:", err)
		return
	}

	response := TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

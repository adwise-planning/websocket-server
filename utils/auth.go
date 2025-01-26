package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
	"websocket-server/database"
	"websocket-server/logging"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("c4e726b73a7b4c91b7e781c6a18e8d2e97fb3f769d47dced8e8d8131a8b4f4a6") // Replace with env variable in production
var tokenExpirationTime = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))               // Token expiration time

// Claims represents the JWT claims
type Claims struct {
	UserName string `json:"username"`
	DeviceID string `json:"device_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new access token
func GenerateToken(username, device_id, email string) (string, error) {
	claims := &Claims{
		UserName: username,
		DeviceID: device_id,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: tokenExpirationTime, // Token expires in 24 hours
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	logging.LogInfo("Token generated successfully for user: " + username)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT and extracts the user ID if valid
func ValidateToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		logging.LogError("Failed to validate token", err)
		// return "", err
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		logging.LogError("Failed to validate token claims", errors.New("invalid token claims"))
		// return "", errors.New("invalid token claims")
	}

	// Check token expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		logging.LogError("Token Expired", errors.New("token has expired"))
		// return "", errors.New("token has expired")
	}

	// Return the user ID from the claims
	logging.LogInfo("Token validated successfully for user: " + claims.UserName)
	return claims.UserName, nil
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken() (string, error) {
	return GenerateRandomString(32) // Implement GenerateRandomString securely
}

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length: must be greater than 0")
	}

	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Encode to a base64 string to ensure safe and readable output
	return base64.RawURLEncoding.EncodeToString(randomBytes)[:length], nil
}

func IsTrustedOrigin(origin string) bool {
	// Implement your logic to check if the origin is trusted
	trustedOrigins := []string{"http://example.com", "http://trusted.com"}
	for _, trustedOrigin := range trustedOrigins {
		if origin == trustedOrigin {
			return true
		}
	}
	return true
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func logToDB(level, message string) {
	query := `INSERT INTO logs (level, message, timestamp) VALUES ($1, $2, $3)`
	_, err := database.PostgresDB.Exec(query, level, message, time.Now())
	if err != nil {
		log.Printf("Failed to log to database: %v", err)
	}
}

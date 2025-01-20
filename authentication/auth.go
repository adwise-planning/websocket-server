package authentication

import (
	"database/sql"
	"errors"
	"time"
	"websocket-server/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("c4e726b73a7b4c91b7e781c6a18e8d2e97fb3f769d47dced8e8d8131a8b4f4a6") // Replace with env variable in production

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new access token
func GenerateJWT(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken() (string, error) {
	return utils.GenerateRandomString(32) // Implement GenerateRandomString securely
}

// GenerateToken generates a JWT for a user
func GenerateToken(userID string) (string, error) {
	// Define the token expiration time
	expirationTime := time.Now().Add(15 * time.Minute) // 15 minutes for access tokens

	// Create the claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
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
		return "", errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	// Check token expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", errors.New("token has expired")
	}

	// Return the user ID from the claims
	return claims.UserID, nil
}

// AuthenticateUser verifies credentials and handles tokens
func AuthenticateUser(db *sql.DB, username, password string) (string, string, error) {
	var userID string
	var hashedPassword string

	query := `SELECT id, password_hash FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		return "", "", errors.New("user not found")
	} else if err != nil {
		return "", "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid password")
	}

	// Generate tokens
	accessToken, err := GenerateJWT(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Store refresh token in DB
	err = StoreRefreshToken(db, userID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// StoreRefreshToken saves a refresh token in the database
func StoreRefreshToken(db *sql.DB, userID, refreshToken string) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, userID, refreshToken, time.Now().Add(30*24*time.Hour)) // 30 days
	return err
}

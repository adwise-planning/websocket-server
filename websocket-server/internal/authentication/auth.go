package authentication

import (
	"database/sql"
	"errors"
	"time"
	"websocket-server/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("c4e726b73a7b4c91b7e781c6a18e8d2e97fb3f769d47dced8e8d8131a8b4f4a6") // Replace with env variable in production

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	return utils.GenerateRandomString(32)
}

func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", errors.New("token has expired")
	}

	return claims.UserID, nil
}

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

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid password")
	}

	accessToken, err := GenerateJWT(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = StoreRefreshToken(db, userID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func StoreRefreshToken(db *sql.DB, userID, refreshToken string) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, userID, refreshToken, time.Now().Add(30*24*time.Hour))
	return err
}
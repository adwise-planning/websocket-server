package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"websocket-server/logging"
	"websocket-server/models"
	"websocket-server/services"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.LogError("Failed to read request body", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for further use

	var requestData models.SignUpRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logging.LogError("Invalid request payload", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := requestData.User
	device := requestData.Device

	s := services.NewUserService()
	// Check if user already exists
	exists := s.UserExists(user.Email)
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Save the user (password should be hashed in production)
	token, refresh_token, err := s.RegisterUser(&user, &device)
	if err != nil {
		logging.LogError("Failed to register user", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	logging.LogInfo(fmt.Sprintf("User registered successfully: %s", user.Username))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "User registered successfully",
		"access_token":  token,
		"refresh_token": refresh_token,
	})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.LogError("Failed to read request body", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for further use

	var requestData models.LoginRequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logging.LogError("Invalid request payload", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	credentials := requestData.Credentials
	device := requestData.Device

	s := services.NewUserService()

	// Validate user credentials
	token, refresh_token, err := s.AuthenticateUser(&credentials, &device)
	if err != nil {
		logging.LogError(fmt.Sprintf("Invalid email or password for user %s", credentials.Username), err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	logging.LogInfo(fmt.Sprintf("User logged in successfully: %s", credentials.Username))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  token,
		"refresh_token": refresh_token,
	})
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"websocket-server/models"
	"websocket-server/services"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		User   models.User   `json:"user"`
		Device models.Device `json:"device"`
	}

	fmt.Printf("Headers: %v\n", r.Header)
	fmt.Printf("Params: %v\n", r.URL.Query())
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Body: %s\n", body)
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for further use

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
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
	log.Printf("%s", fmt.Sprintf("Failed to register user: %v", err))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
		return
	}

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

	var requestData struct {
		Credentials models.Credentials `json:"credentials"`
		Device      models.Device      `json:"device"`
	}

	fmt.Printf("Body: %v\n", r.Body)
	fmt.Printf("JSON Decoder: %v\n", json.NewDecoder(r.Body))

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	credentials := requestData.Credentials
	device := requestData.Device

	s := services.NewUserService()

	// Validate user credentials
	token, refresh_token, err := s.AuthenticateUser(&credentials, &device)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid email or password %s: %v", credentials.Username, err), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"access_token": token, "refresh_token": refresh_token})
}

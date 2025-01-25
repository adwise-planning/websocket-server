package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"websocket-server/database"
	"websocket-server/models"
	"websocket-server/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserService provides user-related functionalities
type UserService struct{}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{}
}

// RegisterUser saves a new user to the database.
func (s *UserService) RegisterUser(user *models.User, device *models.Device) (string, string, error) {
	// Convert password to hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	// Save user to the database
	_, err = database.PostgresDB.Exec(
		"INSERT INTO data.users (username, first_name, last_name, email, date_of_birth, address_line1, address_line2, city, state, country, zip_code, phone_country_code, phone_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		user.Username, user.FirstName, user.LastName, user.Email, user.DateOfBirth, user.AddressLine1, user.AddressLine2, user.City, user.State, user.Country, user.ZipCode, user.PhoneCountryCode, user.PhoneNumber,
	)
	if err != nil {
		return "", "", fmt.Errorf("could not save user: %v", err)
	}

	// Retrieve the user ID and save user_auth
	var userID int
	err = database.PostgresDB.QueryRow("SELECT user_id FROM data.users WHERE email=$1", user.Email).Scan(&userID)
	if err != nil {
		return "", "", fmt.Errorf("could not retrieve user ID: %v", err)
	}

	// Save device to the database
	s.SaveDevice(userID, device)

	// Generate auth token
	accessToken, err := utils.GenerateToken(user.Username, device.DeviceID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate token: %v", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Save user_auth to the database
	_, err = database.PostgresDB.Exec(
		"INSERT INTO data.user_auth (user_id, password_hash, auth_token, refresh_token) VALUES ($1, $2, $3, $4)",
		userID, string(hashedPassword), accessToken, refreshToken,
	)
	if err != nil {
		return "", "", fmt.Errorf("could not save user authentication: %v", err)
	}
	return accessToken, refreshToken, nil
}

// GetToken retrieves the token for a user from the database - Login User
func (s *UserService) AuthenticateUser(credentials *models.Credentials, device *models.Device) (string, string, error) {
	var user_id int
	var password_hash, email string
	query := "SELECT a.auth_id, password_hash, u.email FROM data.user_auth a join data.users u on a.user_id = u.user_id WHERE u.username=$1"
	err := database.PostgresDB.QueryRow(query, credentials.Username).Scan(&user_id, &password_hash, &email)
	if err == sql.ErrNoRows {
		return "", "", errors.New("user not found")
	} else if err != nil {
		return "", "", fmt.Errorf("error for %s: %v", credentials.Username, err)
	}

	// Compare the stored password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(credentials.Password))
	if err != nil {
		return "", "", errors.New("invalid Credentials")
	}

	accessToken, err := utils.GenerateToken(credentials.Username, device.DeviceID, email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Save the device
	s.SaveDevice(user_id, device)

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Store refresh token in DB
	err = s.StoreRefreshToken(user_id, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// StoreRefreshToken saves a refresh token in the database
func (s *UserService) StoreRefreshToken(userID int, refreshToken string) error {
	query := `INSERT INTO data.refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := database.PostgresDB.Exec(query, userID, refreshToken, time.Now().Add(30*24*time.Hour)) // 30 days
	return err
}

// SaveDevice saves a new device to the database.
func (s *UserService) SaveDevice(user_id int, device *models.Device) (string, error) {
	_, err := database.PostgresDB.Exec(
		`INSERT INTO data.devices (
			user_id, name, type, manufacturer, model, serial_number, firmware, hardware_version, 
			software_version, operating_system, processor, memory, storage_capacity, screen_size, 
			resolution, camera, sensors, ports, dimensions, weight, color, material, power_source, 
			battery_level, signal_strength, connectivity_type, ip_address, mac_address, network_provider, 
			plan_type, subscription_end, status, last_seen, location, owner, created_at, updated_at, 
			usage_time, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
		)`,
		user_id, device.Name, device.Type, device.Manufacturer, device.Model, device.SerialNumber,
		device.Firmware, device.HardwareVersion, device.SoftwareVersion, device.OperatingSystem,
		device.Processor, device.Memory, device.StorageCapacity, device.ScreenSize, device.Resolution,
		device.Camera, device.Sensors, device.Ports, device.Dimensions, device.Weight, device.Color,
		device.Material, device.PowerSource, device.BatteryLevel, device.SignalStrength,
		device.ConnectivityType, device.IPAddress, device.MACAddress, device.NetworkProvider,
		device.PlanType, device.SubscriptionEnd, device.Status, device.LastSeen, device.Location,
		device.Owner, device.CreatedAt, device.UpdatedAt, device.UsageTime, device.Notes,
	)
	if err != nil {
		return "", fmt.Errorf("could not save device: %v", err)
	}
	return "Device saved successfully", nil
}

// UserExists checks if a user with the given email already exists in the database.
func (s *UserService) UserExists(email string) bool {
	var exists bool = true
	query := "SELECT EXISTS (SELECT 1 FROM data.users WHERE email=$1)"
	err := database.PostgresDB.QueryRow(query, email).Scan(&exists)
	log.Printf("could not check if user exists: %v", err)
	if err != nil {
		return false
	}
	return exists
}

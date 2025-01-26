package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"websocket-server/database"
	"websocket-server/logging"
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
		logging.LogError("Failed to hash password", err)
		return "", "", err
	}

	// Save user to the database
	tx, err := database.PostgresDB.Begin()
	if err != nil {
		logging.LogError("Failed to begin transaction", err)
		return "", "", err
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRow(
		"INSERT INTO data.users (username, first_name, last_name, email, date_of_birth, address_line1, address_line2, city, state, country, zip_code, phone_country_code, phone_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING user_id",
		user.Username, user.FirstName, user.LastName, user.Email, user.DateOfBirth, user.AddressLine1, user.AddressLine2, user.City, user.State, user.Country, user.ZipCode, user.PhoneCountryCode, user.PhoneNumber,
	).Scan(&userID)
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to save user: %s", user.Username), err)
		return "", "", fmt.Errorf("could not save user: %v", err)
	}
	if err != nil {
		logging.LogError(fmt.Sprintf("Could not retrieve user ID for username: %s", user.Username), err)
		return "", "", fmt.Errorf("could not retrieve user ID: %v", err)
	}
	logging.LogInfo("Device saved successfully for user ID: " + fmt.Sprint(userID))

	// Save device to the database
	deviceID, err := s.SaveDevice(tx, userID, device)
	if err != nil {
		return "", "", err
	}

	// Generate auth token
	accessToken, err := utils.GenerateToken(user.Username, deviceID, user.Email)
	if err != nil {
		logging.LogError("Failed to generate access token", err)
		return "", "", fmt.Errorf("could not generate token: %v", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		logging.LogError("Failed to generate refresh token", err)
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Save user_auth to the database
	_, err = tx.Exec(
		"INSERT INTO data.user_auth (user_id, password_hash, auth_token, refresh_token) VALUES ($1, $2, $3, $4)",
		userID, string(hashedPassword), accessToken, refreshToken,
	)
	if err != nil {
		logging.LogError(fmt.Sprintf("Could not save user authentication details in DB for user ID: %d", userID), err)
		return "", "", fmt.Errorf("could not save user authentication details in DB: %v", err)
	}

	if err = tx.Commit(); err != nil {
		logging.LogError("Failed to commit transaction", err)
		return "", "", err
	}

	logging.LogInfo("User registered successfully: " + user.Username)
	return accessToken, refreshToken, nil
}

// GetToken retrieves the token for a user from the database - Login User
func (s *UserService) AuthenticateUser(credentials *models.Credentials, device *models.Device) (string, string, error) {
	var userID int64
	var password_hash, email string
	query := "SELECT a.auth_id, password_hash, u.email FROM data.user_auth a join data.users u on a.user_id = u.user_id WHERE u.username=$1"
	err := database.PostgresDB.QueryRow(query, credentials.Username).Scan(&userID, &password_hash, &email)
	if err == sql.ErrNoRows {
		logging.LogError("User not found", err)
		return "", "", errors.New("user not found")
	} else if err != nil {
		logging.LogError("Error querying user", err)
		return "", "", fmt.Errorf("error querying user: %v", err)
	}

	// Compare the stored password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(credentials.Password))
	if err != nil {
		logging.LogError("Invalid credentials", err)
		return "", "", errors.New("invalid Credentials")
	}

	accessToken, err := utils.GenerateToken(credentials.Username, device.DeviceID, email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Save the device
	s.SaveDevice(nil, userID, device)

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		logging.LogError("Failed to generate refresh token", err)
		return "", "", fmt.Errorf("could not generate access token: %v", err)
	}

	// Store refresh token in DB
	err = s.StoreRefreshToken(userID, refreshToken)
	if err != nil {
		logging.LogError("Failed to store refresh token", err)
		return "", "", err
	}

	logging.LogInfo(fmt.Sprintf("User authenticated successfully: %s", credentials.Username))
	return accessToken, refreshToken, nil
}

// StoreRefreshToken saves a refresh token in the database
func (s *UserService) StoreRefreshToken(userID int64, refreshToken string) error {
	query := `INSERT INTO data.refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := database.PostgresDB.Exec(query, userID, refreshToken, time.Now().Add(30*24*time.Hour)) // 30 days
	if err != nil {
		logging.LogError("Failed to store refresh token", err)
		return nil
		return err
	}
	logging.LogInfo(fmt.Sprintf("Refresh token stored successfully for user ID: %d", userID))
	return nil
}

// SaveDevice saves a new device to the database.
func (s *UserService) SaveDevice(tx *sql.Tx, userID int64, device *models.Device) (string, error) {
	query := `INSERT INTO data.devices (
			user_id, name, type, manufacturer, model, serial_number, firmware, hardware_version, 
			software_version, operating_system, processor, memory, storage_capacity, screen_size, 
			resolution, camera, sensors, ports, dimensions, weight, color, material, power_source, 
			battery_level, signal_strength, connectivity_type, ip_address, mac_address, network_provider, 
			plan_type, subscription_end, status, last_seen, location, owner, created_at, updated_at, 
			usage_time, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
		) RETURNING device_id`
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.Exec(query, userID, device.Name, device.Type, device.Manufacturer, device.Model, device.SerialNumber,
			device.Firmware, device.HardwareVersion, device.SoftwareVersion, device.OperatingSystem,
			device.Processor, device.Memory, device.StorageCapacity, device.ScreenSize, device.Resolution,
			device.Camera, device.Sensors, device.Ports, device.Dimensions, device.Weight, device.Color,
			device.Material, device.PowerSource, device.BatteryLevel, device.SignalStrength,
			device.ConnectivityType, device.IPAddress, device.MACAddress, device.NetworkProvider,
			device.PlanType, device.SubscriptionEnd, device.Status, device.LastSeen, device.Location,
			device.Owner, device.CreatedAt, device.UpdatedAt, device.UsageTime, device.Notes)
	} else {
		result, err = database.PostgresDB.Exec(query, userID, device.Name, device.Type, device.Manufacturer, device.Model, device.SerialNumber,
			device.Firmware, device.HardwareVersion, device.SoftwareVersion, device.OperatingSystem,
			device.Processor, device.Memory, device.StorageCapacity, device.ScreenSize, device.Resolution,
			device.Camera, device.Sensors, device.Ports, device.Dimensions, device.Weight, device.Color,
			device.Material, device.PowerSource, device.BatteryLevel, device.SignalStrength,
			device.ConnectivityType, device.IPAddress, device.MACAddress, device.NetworkProvider,
			device.PlanType, device.SubscriptionEnd, device.Status, device.LastSeen, device.Location,
			device.Owner, device.CreatedAt, device.UpdatedAt, device.UsageTime, device.Notes)
	}

	if err != nil {
		logging.LogError("Failed to save device to database", err)
		return "", fmt.Errorf("could not save device: %v", err)
	}
	deviceID, err := result.LastInsertId()
	if err != nil {
		logging.LogError("Failed to retrieve device id from database", err)
		return "", fmt.Errorf("could not retrieve device ID: %v", err)
	}
	logging.LogInfo("Device saved successfully for user ID: " + fmt.Sprint(userID))
	return fmt.Sprintf("Device %d saved successfully", deviceID), nil
}

// UserExists checks if a user with the given email already exists in the database.
func (s *UserService) UserExists(email string) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM data.users WHERE email=$1)"
	err := database.PostgresDB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		logging.LogError("Could not check if user exists", err)
		return false
	}
	return exists
}

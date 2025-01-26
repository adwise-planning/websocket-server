package models

import (
	"errors"
	"regexp"
	"time"
)

// Data structure for user registration request
type SignUpRequestData struct {
	User   User   `json:"user"`
	Device Device `json:"device"`
}

// Data structure for user login request
type LoginRequestData struct {
	Credentials Credentials `json:"credentials"`
	Device      Device      `json:"device"`
}

// Credentials for login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User represents a user in the system
type User struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	DateOfBirth      string `json:"date_of_birth"`
	AddressLine1     string `json:"address_line_1"`
	AddressLine2     string `json:"address_line_2"`
	City             string `json:"city"`
	State            string `json:"state"`
	Country          string `json:"country"`
	ZipCode          string `json:"zip_code"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
}

// User Device details for logging
type Device struct {
	DeviceID         string `json:"device_id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Manufacturer     string `json:"manufacturer"`
	Model            string `json:"model"`
	SerialNumber     string `json:"serial_number"`
	Firmware         string `json:"firmware"`
	HardwareVersion  string `json:"hardware_version"`
	SoftwareVersion  string `json:"software_version"`
	OperatingSystem  string `json:"operating_system"`
	Processor        string `json:"processor"`
	Memory           int    `json:"memory"`
	StorageCapacity  int    `json:"storage_capacity"`
	ScreenSize       string `json:"screen_size"`
	Resolution       string `json:"resolution"`
	Camera           string `json:"camera"`
	Sensors          string `json:"sensors"`
	Ports            string `json:"ports"`
	Dimensions       string `json:"dimensions"`
	Weight           string `json:"weight"`
	Color            string `json:"color"`
	Material         string `json:"material"`
	PowerSource      string `json:"power_source"`
	BatteryLevel     int    `json:"battery_level"`
	SignalStrength   int    `json:"signal_strength"`
	ConnectivityType string `json:"connectivity_type"`
	IPAddress        string `json:"ip_address"`
	MACAddress       string `json:"mac_address"`
	NetworkProvider  string `json:"network_provider"`
	PlanType         string `json:"plan_type"`
	SubscriptionEnd  string `json:"subscription_end"`
	Status           string `json:"status"`
	LastSeen         string `json:"last_seen"`
	Location         string `json:"location"`
	Owner            string `json:"owner"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	UsageTime        string `json:"usage_time"`
	Notes            string `json:"notes"`
}

// Validate validates the User struct fields
func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" || u.Email == "" {
		return errors.New("username, password, and email are required")
	}
	if len(u.Username) < 3 || len(u.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.DateOfBirth != "" {
		if _, err := time.Parse("2006-01-02", u.DateOfBirth); err != nil {
			return errors.New("invalid date of birth format, expected YYYY-MM-DD")
		}
	}
	if u.PhoneNumber != "" && !isValidPhoneNumber(u.PhoneNumber) {
		return errors.New("invalid phone number format")
	}
	if u.ZipCode != "" {
		re := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
		if !re.MatchString(u.ZipCode) {
			return errors.New("invalid zip code format")
		}
	}
	if u.Country == "" {
		return errors.New("country is required")
	}
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("first name and last name are required")
	}
	if u.AddressLine1 == "" {
		return errors.New("address line 1 is required")
	}
	if u.City == "" {
		return errors.New("city is required")
	}
	if u.State == "" {
		return errors.New("state is required")
	}
	if u.PhoneCountryCode == "" {
		return errors.New("phone country code is required")
	}
	if u.Password == u.Username {
		return errors.New("password cannot be the same as username")
	}
	if u.Password == u.Email {
		return errors.New("password cannot be the same as email")
	}
	if u.Password == u.FirstName || u.Password == u.LastName {
		return errors.New("password cannot be the same as first name or last name")
	}
	if !containsUppercase(u.Password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !containsLowercase(u.Password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !containsDigit(u.Password) {
		return errors.New("password must contain at least one digit")
	}
	if !containsSpecialChar(u.Password) {
		return errors.New("password must contain at least one special character")
	}
	if containsWhitespace(u.Password) {
		return errors.New("password cannot contain whitespace")
	}
	if len(u.Password) > 64 {
		return errors.New("password cannot be longer than 64 characters")
	}
	if isCommonPassword(u.Password) {
		return errors.New("password is too common")
	}
	return nil
}

func containsUppercase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func containsSpecialChar(s string) bool {
	re := regexp.MustCompile(`[!@#~$%^&*()_+|<>?:{}]`)
	return re.MatchString(s)
}

func containsWhitespace(s string) bool {
	for _, c := range s {
		if c == ' ' {
			return true
		}
	}
	return false
}

func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"123456", "password", "123456789", "12345678", "12345", "1234567", "1234567890", "qwerty", "abc123", "password1",
		"111111", "123123", "admin", "letmein", "welcome", "monkey", "dragon", "football", "iloveyou", "sunshine",
		"princess", "admin123", "1234", "passw0rd", "master", "hello", "freedom", "whatever", "qazwsx", "trustno1",
	}
	for _, p := range commonPasswords {
		if password == p {
			return true
		}
	}
	return false
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// isValidPhoneNumber checks if the phone number format is valid
func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

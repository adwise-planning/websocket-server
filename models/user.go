package models

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

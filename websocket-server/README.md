# websocket-server

## Overview
The websocket-server project is a backend service that provides enterprise-level WebSocket-based messaging functionality. It includes user authentication, message logging, security logging, and user event tracking, all stored in a PostgreSQL database.

## Features
- WebSocket-based messaging for real-time communication.
- User authentication using JWT.
- RESTful API for user management and message handling.
- Logging of user actions and messages for security and auditing.
- PostgreSQL database for persistent storage.

## Project Structure
```
websocket-server
├── cmd
│   └── main.go                # Entry point of the application
├── config
│   └── config.go              # Configuration settings
├── internal
│   ├── authentication          # User authentication logic
│   │   └── auth.go
│   ├── handlers                # HTTP and WebSocket handlers
│   │   ├── websocket.go
│   │   └── api.go
│   ├── logging                 # Logging functionality
│   │   └── logger.go
│   ├── models                  # Data models
│   │   ├── user.go
│   │   ├── message.go
│   │   └── event.go
│   ├── repository              # Database interaction
│   │   ├── user_repository.go
│   │   ├── message_repository.go
│   │   └── event_repository.go
│   ├── services                # Business logic
│   │   ├── user_service.go
│   │   ├── message_service.go
│   │   └── event_service.go
│   └── utils                   # Utility functions
│       └── utils.go
├── migrations                   # Database migrations
│   └── 001_create_tables.sql
├── Dockerfile                   # Docker configuration
├── go.mod                       # Module dependencies
├── go.sum                       # Dependency checksums
└── README.md                    # Project documentation
```

## Setup Instructions
1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd websocket-server
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Configure the database:**
   Update the `config/config.go` file with your PostgreSQL connection details.

4. **Run database migrations:**
   Execute the SQL commands in `migrations/001_create_tables.sql` to set up the necessary tables.

5. **Start the application:**
   ```
   go run cmd/main.go
   ```

## Usage
- Connect to the WebSocket server for real-time messaging.
- Use the RESTful API for user management and message handling.

## Logging
All user actions, messages, and events are logged for security and auditing purposes. Check the logs for detailed information on user interactions.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.
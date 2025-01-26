package logging

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
)

// Initialize the logger
var (
	infoLogger        = newAsyncLogger(os.Stdout, "INFO: ")
	securityLogger    = newAsyncLogger(os.Stdout, "SECURITY: ")
	eventLogger       = newAsyncLogger(os.Stdout, "EVENT: ")
	errorLogger       = newAsyncLogger(os.Stderr, "ERROR: ")
	debugLogger       = newAsyncLogger(os.Stdout, "DEBUG: ")
	warnLogger        = newAsyncLogger(os.Stdout, "WARN: ")
	accessLogger      = newAsyncLogger(os.Stdout, "ACCESS: ")
	performanceLogger = newAsyncLogger(os.Stdout, "PERFORMANCE: ")
	fatalLogger       = newAsyncLogger(os.Stderr, "FATAL: ")
)

// asyncLogger wraps a standard logger to provide asynchronous logging
type asyncLogger struct {
	logger *log.Logger
	mu     sync.Mutex
}

// newAsyncLogger creates a new asyncLogger
func newAsyncLogger(out *os.File, prefix string) *asyncLogger {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
	}
	return &asyncLogger{
		logger: log.New(out, prefix+"["+file+"] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Printf logs a formatted message asynchronously
func (a *asyncLogger) Printf(format string, v ...interface{}) {
	go func() {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.logger.Printf(format, v...)
	}()
}

// Println logs a message asynchronously
func (a *asyncLogger) Println(v ...interface{}) {
	go func() {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.logger.Println(v...)
	}()
}

// func logToDB(level, message string) {
// 	query := `INSERT INTO logs (level, message, timestamp) VALUES ($1, $2, $3)`
// 	_, err := database.PostgresDB.Exec(query, level, message, time.Now())
// 	if err != nil {
// 		log.Printf("Failed to log to database: %v", err)
// 	}
// }

// LogSecurityEvent logs security-related events
func LogSecurityEvent(message, username string) {
	securityLogger.Printf("%s - User: %s", message, username)
	// logToDB("SECURITY", message+" - User: "+username)
}

// LogEvent logs general events
func LogEvent(message string) {
	eventLogger.Println(message)
	// logToDB("EVENT", message)
}

// LogError logs errors with context
func LogError(message string, err error) {
	errorLogger.Printf("%s: %v", message, err)
	// logToDB("ERROR", message+": "+err.Error())
}

// LogInfo logs informational messages
func LogInfo(message string) {
	infoLogger.Println(message)
	// logToDB("INFO", message)
}

// LogDebug logs debug messages
func LogDebug(message string) {
	debugLogger.Println(message)
	// logToDB("DEBUG", message)
}

// LogWarn logs warning messages
func LogWarn(message string) {
	warnLogger.Println(message)
	// logToDB("WARN", message)
}

// LogAccess logs access-related events
func LogAccess(message string) {
	accessLogger.Println(message)
	// logToDB("ACCESS", message)
}

// LogPerformance logs performance-related events
func LogPerformance(message string) {
	performanceLogger.Println(message)
	// logToDB("PERFORMANCE", message)
}

// LogFatal logs fatal errors and exits the application
func LogFatal(message string) {
	fatalLogger.Println(message)
	// logToDB("FATAL", message)
}

// LogErrorAndRespond logs the error and sends an HTTP response with the given status code and message.
func LogErrorAndRespond(w http.ResponseWriter, message string, statusCode int, err error) {
	LogError(message, err)
	http.Error(w, message, statusCode)
}

package logging

import (
	"log"
	"os"
)

var (
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	securityLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	securityLogger = log.New(file, "SECURITY: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	infoLogger.Println(message)
}

func Error(message string) {
	errorLogger.Println(message)
}

func Security(message string) {
	securityLogger.Println(message)
}
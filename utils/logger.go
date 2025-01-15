package utils

import (
	"log"
	"os"
)

// Logger wraps standard log functions
var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

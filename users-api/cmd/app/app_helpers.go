package app

import (
	"log"
	"os"
)

// sanityCheck checks that all required environment variables are set.
// if any of the required variables is not defined, it prints a log message.
func sanityCheck(l *log.Logger) {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			l.Printf("environment variable %s not defined. Terminating application...", k)
		}
	}
}

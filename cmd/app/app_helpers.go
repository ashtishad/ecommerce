package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// sanityCheck checks that all required environment variables are set.
// if any of the required variables is not defined, it prints a log message.
func sanityCheck() {
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
			log.Println(fmt.Sprintf("environment variable %s not defined. Terminating application...", k))
		}
	}
}

// writeResponse writes api endpoint response data and correct http status code in response.
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

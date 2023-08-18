package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func Start() {
	sanityCheck()

	r := mux.NewRouter()

	// database connection config
	conn := getDbClient()
	defer conn.Close()

	// define routes
	r.
		HandleFunc("/user", createUserHandler).
		Methods(http.MethodPost).
		Name("Create User")

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

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
			log.Println(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}

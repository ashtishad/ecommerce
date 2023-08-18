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

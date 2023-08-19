package app

import (
	"bitbucket.org/ashtishad/as_ti/domain"
	"bitbucket.org/ashtishad/as_ti/service"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func Start() {
	sanityCheck()

	r := mux.NewRouter()

	// initiated logger, dependency injection, create once, inject it where needed
	l := log.New(os.Stdout, "users-api ", log.LstdFlags)

	// database connection config
	conn := getDbClient()
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			l.Printf("couldn't close the database client : %v", err.Error())
		}
	}(conn)

	// wire up the handler
	userRepositoryDB := domain.NewUserRepositoryDB(conn, l)
	uh := UserHandlers{service.NewUserService(userRepositoryDB), l}

	// define routes
	r.
		HandleFunc("/user", uh.createUserHandler).
		Methods(http.MethodPost).
		Name("Create User")

	r.
		HandleFunc("/existing-user", uh.existingUserHandler).
		Methods(http.MethodPost).
		Name("Existing User")

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

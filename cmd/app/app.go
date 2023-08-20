package app

import (
	"database/sql"
	"fmt"
	"github.com/ashtishad/ecommerce/domain"
	"github.com/ashtishad/ecommerce/service"
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
	gh := GoogleAuthHandler{l: l}

	// define routes
	r.
		HandleFunc("/user", uh.createUserHandler).
		Methods(http.MethodPost).
		Name("Create User")

	r.
		HandleFunc("/existing-user", uh.existingUserHandler).
		Methods(http.MethodPost, http.MethodGet).
		Name("Existing User")
	r.
		HandleFunc("/login", gh.startGoogleLoginHandler).
		Methods(http.MethodPost, http.MethodGet).
		Name("Google Login")
	r.
		HandleFunc("/callback", gh.googleCallbackHandler).
		Methods(http.MethodPost, http.MethodGet).
		Name("Google Callback")

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

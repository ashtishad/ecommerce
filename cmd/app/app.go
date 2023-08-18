package app

import (
	"bitbucket.org/ashtishad/as_ti/domain"
	"bitbucket.org/ashtishad/as_ti/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func Start() {
	sanityCheck()

	// l := log.New(os.Stdout, "users-api ", log.LstdFlags)

	r := mux.NewRouter()

	// database connection config
	conn := getDbClient()
	defer conn.Close()

	// wire up the handler
	userRepositoryDB := domain.NewUserRepositoryDB(conn)
	uh := UserHandlers{service.NewUserService(userRepositoryDB)}
	// define routes
	r.
		HandleFunc("/user", uh.createUserHandler).
		Methods(http.MethodPost).
		Name("Create User")

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

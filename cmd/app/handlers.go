package app

import (
	"fmt"
	"net/http"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Successfully handled request for create user")
}

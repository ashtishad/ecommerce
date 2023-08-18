package app

import (
	"bitbucket.org/ashtishad/as_ti/domain"
	"bitbucket.org/ashtishad/as_ti/service"
	"encoding/json"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
}

func (us *UserHandlers) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request
	var newUserRequest domain.NewUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&newUserRequest)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the request data as needed using regex

	// Call the service method to create a new user
	userResponse, err := us.service.NewUser(newUserRequest)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

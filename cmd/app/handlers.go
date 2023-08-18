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

// createUserHandler decodes the user request, returns bad request error if failed to decode json
// then validates user data using regex,
// then calls the service method to create a new user,
// finally write the response data and correct http status code.
func (us *UserHandlers) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUserRequest domain.NewUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&newUserRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}

	// Validate the request data as needed using regex
	// ...

	userResponse, err := us.service.NewUser(newUserRequest)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeResponse(w, http.StatusOK, userResponse)
}

package app

import (
	"bitbucket.org/ashtishad/as_ti/domain"
	"bitbucket.org/ashtishad/as_ti/service"
	"encoding/json"
	"log"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
	l       *log.Logger
}

// createUserHandler decodes the user request, returns bad request error if failed to decode json
// then validates user data using regex,
// then calls the service method to create a new user,
// finally write the response data and correct http status code.
func (us *UserHandlers) createUserHandler(w http.ResponseWriter, r *http.Request) {
	us.l.Println("Handling POST request on /user")

	var newUserRequest domain.NewUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&newUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}

	if err := validateCreateUserInput(newUserRequest); err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userResponse, err := us.service.NewUser(newUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeResponse(w, http.StatusOK, userResponse)
}

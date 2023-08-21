package app

import (
	"encoding/json"
	"github.com/ashtishad/ecommerce/domain"
	"github.com/ashtishad/ecommerce/service"
	"github.com/gorilla/mux"
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
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeResponse(w, http.StatusOK, userResponse)
}

// updateUserHandler decodes the user request, returns bad request error if failed to decode json
// then validates user data using regex,
// then calls the service method to create a new user,
// finally write the response data and correct http status code.
func (us *UserHandlers) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	us.l.Println("Handling POST request on /user")

	vars := mux.Vars(r)
	UserUUID := vars["user_id"]

	var updateUserRequest domain.UpdateUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}

	updateUserRequest.UserUUID = UserUUID
	if err := validateUpdateUserInput(updateUserRequest); err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userResponse, err := us.service.UpdateUser(updateUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeResponse(w, http.StatusOK, userResponse)
}

// existingUserHandler decodes the existing user request, returns bad request error if failed to decode json.
// then validates user's email and password data using regex,
// then calls the service method to get existing user,
// finally write the response data and correct http status code
func (us *UserHandlers) existingUserHandler(w http.ResponseWriter, r *http.Request) {
	us.l.Println("Handling GET request on /user")

	var existingUserRequest domain.ExistingUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&existingUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validateExistingUserInput(existingUserRequest); err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userResponse, err := us.service.ExistingUser(existingUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeResponse(w, http.StatusOK, userResponse)
}

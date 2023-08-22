package app

import (
	"github.com/ashtishad/ecommerce/domain"
	"github.com/ashtishad/ecommerce/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
	l       *log.Logger
}

// createUserHandler handles the creation of a user if user not exists.
// It decodes the user request, returns bad request error if failed to decode json,
// then calls the service method to create a new user,
// finally write the response data and correct HTTP status code.
func (us *UserHandlers) createUserHandler(c *gin.Context) {
	us.l.Println("Handling POST request on /user")

	var newUserRequest domain.NewUserRequestDTO
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	userResponse, err := us.service.NewUser(newUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// updateUserHandler handles the updating of a user.
// It decodes the user request, returns bad request error if failed to decode json,
// then calls the service method to update the user,
// finally write the response data and correct HTTP status code.
func (us *UserHandlers) updateUserHandler(c *gin.Context) {
	us.l.Println("Handling PUT request on /user")

	UserUUID := c.Param("user_id")

	var updateUserRequest domain.UpdateUserRequestDTO
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	updateUserRequest.UserUUID = UserUUID
	userResponse, err := us.service.UpdateUser(updateUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// existingUserHandler handles the retrieval of an existing user.
// It decodes the existing user request, returns bad request error if failed to decode json,
// then calls the service method to get the existing user,
// finally write the response data and correct HTTP status code.
func (us *UserHandlers) existingUserHandler(c *gin.Context) {
	us.l.Println("Handling GET request on /user")

	var existingUserRequest domain.ExistingUserRequestDTO
	if err := c.ShouldBindJSON(&existingUserRequest); err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := us.service.ExistingUser(existingUserRequest)
	if err != nil {
		us.l.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

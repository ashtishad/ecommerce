package app

import (
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
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

// GetUsersHandler handles the GET request to fetch all users
func (us *UserHandlers) GetUsersHandler(c *gin.Context) {
	var opts domain.FindAllUsersOptionsDTO
	opts.FromIDStr = c.Query("fromID")
	opts.PageSizeStr = c.DefaultQuery("pageSize", constants.DefaultPageSizeString)
	opts.Status = c.DefaultQuery("status", constants.UserStatusActive)
	opts.SignUpOption = c.Query("signUpOption")
	opts.Timezone = c.Query("timezone")

	users, pageInfo, err := us.service.GetAllUsers(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"page_info": pageInfo,
	})
}

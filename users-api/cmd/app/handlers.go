package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	service service.UserService
	l       *slog.Logger
}

// createUserHandler handles the creation of a user if user not exists.
// It decodes the user request, returns bad request error if failed to decode json,
// then calls the service method to create a new user,
// finally write the response data and correct HTTP status code.
func (us *UserHandlers) createUserHandler(c *gin.Context) {
	var newUserRequest domain.NewUserRequestDTO
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		us.l.Error("failed to bind create user req dto", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})

		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), lib.TimeoutCreateUser)
	defer cancel()

	userResponse, err := us.service.NewUser(timeoutCtx, newUserRequest)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// updateUserHandler handles the updating of a user.
// It decodes the user request, returns bad request error if failed to decode json,
// then calls the service method to update the user,
// finally write the response data and correct HTTP status code.
func (us *UserHandlers) updateUserHandler(c *gin.Context) {
	UserUUID := c.Param("user_id")

	var updateUserRequest domain.UpdateUserRequestDTO
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		us.l.Error("failed to bind update user req dto", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})

		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), lib.TimeoutUpdateUser)
	defer cancel()

	updateUserRequest.UserUUID = UserUUID
	userResponse, err := us.service.UpdateUser(timeoutCtx, updateUserRequest)

	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// GetUsersHandler handles the GET request to fetch all users
func (us *UserHandlers) GetUsersHandler(c *gin.Context) {
	var opts domain.FindAllUsersOptionsDTO
	opts.FromIDStr = c.Query("fromID")
	opts.PageSizeStr = c.Query("pageSize")
	opts.Status = c.Query("status")
	opts.SignUpOption = c.Query("signUpOption")
	opts.Timezone = c.Query("timezone")

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), lib.TimeoutGetUsers)
	defer cancel()

	users, pageInfo, err := us.service.GetAllUsers(timeoutCtx, opts)

	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"page_info": pageInfo,
	})
}

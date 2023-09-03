package app

import (
	"bytes"
	"encoding/json"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	mockUserService := new(service.MockUserService)

	expectedUserResponse := &domain.UserResponseDTO{
		UserUUID:     "some-uuid",
		Email:        "john_wick@gmail.com",
		FullName:     "John Wick",
		Phone:        "1234567890",
		SignUpOption: "general",
		Status:       "active",
		Timezone:     "Asia/Dhaka",
	}

	mockUserService.On("NewUser", mock.AnythingOfType("domain.NewUserRequestDTO")).Return(expectedUserResponse, nil)

	us := UserHandlers{service: mockUserService}

	newUserRequest := domain.NewUserRequestDTO{
		Email:        "john_wick@gmail.com",
		FullName:     "John Wick",
		Password:     "secret_pass",
		Phone:        "1234567890",
		SignUpOption: "general",
		Timezone:     "Asia/Dhaka",
	}

	data, _ := json.Marshal(newUserRequest)

	r.POST("/users", us.createUserHandler)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var receivedUser domain.UserResponseDTO
	err := json.NewDecoder(resp.Body).Decode(&receivedUser)
	assert.Nil(t, err)

	assert.Equal(t, expectedUserResponse, &receivedUser)
}

func TestUpdateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	mockUserService := new(service.MockUserService)

	expectedUserResponse := &domain.UserResponseDTO{
		UserUUID:     "some-uuid",
		Email:        "keanu_reeves@gmail.com",
		FullName:     "Keanu Reeves",
		Phone:        "9876543210",
		SignUpOption: "general",
		Status:       "active",
		Timezone:     "Asia/Dhaka",
	}

	updateUserRequest := domain.UpdateUserRequestDTO{
		UserUUID: "some-uuid",
		Email:    "keanu_reeves@gmail.com",
		FullName: "Keanu Reeves",
		Phone:    "9876543210",
		Timezone: "Asia/Dhaka",
	}

	mockUserService.On("UpdateUser", mock.AnythingOfType("domain.UpdateUserRequestDTO")).Return(expectedUserResponse, nil)

	us := UserHandlers{service: mockUserService}

	data, _ := json.Marshal(updateUserRequest)

	r.PUT("/users/:user_id", us.updateUserHandler)

	req, _ := http.NewRequest("PUT", "/users/some-uuid", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var receivedUser domain.UserResponseDTO
	err := json.NewDecoder(resp.Body).Decode(&receivedUser)
	assert.Nil(t, err)

	assert.Equal(t, expectedUserResponse, &receivedUser)
}

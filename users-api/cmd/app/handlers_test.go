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
		Email:        "keanu_reeves@gmail.com",
		FullName:     "Keanu Reeves",
		Phone:        "1234567890",
		SignUpOption: "general",
		Status:       "active",
		Timezone:     "Asia/Dhaka",
	}

	mockUserService.On("NewUser", mock.AnythingOfType("domain.NewUserRequestDTO")).Return(expectedUserResponse, nil)

	us := UserHandlers{service: mockUserService}

	newUserRequest := domain.NewUserRequestDTO{
		Email:        "KeeanuReeves@outlook.com",
		FullName:     "Keanu Reeves",
		Password:     "secrpsswrd",
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

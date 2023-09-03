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

func TestGetUsersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserService := new(service.MockUserService)
	us := UserHandlers{service: mockUserService}

	expectedUsers := []domain.UserResponseDTO{
		{
			UserUUID:     "uuid1",
			Email:        "user1@example.com",
			FullName:     "User One",
			Phone:        "0198747782808",
			SignUpOption: "general",
			Status:       "active",
			Timezone:     "America/New_York",
		},
		{
			UserUUID:     "uuid2",
			Email:        "user2@example.com",
			FullName:     "User Two",
			Phone:        "0198747782808",
			SignUpOption: "google",
			Status:       "active",
			Timezone:     "Australia/Sydney",
		},
	}

	expectedPageInfo := &domain.NextPageInfo{
		HasNextPage: true,
		StartCursor: 1,
		EndCursor:   2,
		TotalCount:  4,
	}

	mockUserService.On("GetAllUsers", mock.AnythingOfType("domain.FindAllUsersOptionsDTO")).
		Return(expectedUsers, expectedPageInfo, nil)

	r.GET("/users", us.GetUsersHandler)

	testCases := []struct {
		queryString string
	}{
		{""},
		{"?fromID=1&pageSize=2"},
		{"?status=active"},
		{"?signUpOption=general"},
		{"?timezone=Asia/Dhaka"},
		{"?fromID=1&pageSize=2&status=active&signUpOption=general&timezone=Asia/Dhaka"},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/users"+tc.queryString, nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result struct {
			Users    []domain.UserResponseDTO `json:"users"`
			PageInfo *domain.NextPageInfo     `json:"page_info"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.Nil(t, err)

		assert.Equal(t, expectedUsers, result.Users)
		assert.Equal(t, expectedPageInfo, result.PageInfo)
	}
}

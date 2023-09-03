package service

import (
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) NewUser(request domain.NewUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}
	if args.Get(0) != nil {
		return args.Get(0).(*domain.UserResponseDTO), nil
	}
	return nil, args.Get(1).(lib.APIError)
}

func (m *MockUserService) UpdateUser(request domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}
	if args.Get(0) != nil {
		return args.Get(0).(*domain.UserResponseDTO), nil
	}
	return nil, args.Get(1).(lib.APIError)
}

func (m *MockUserService) GetAllUsers(request domain.FindAllUsersOptionsDTO) ([]domain.UserResponseDTO, *domain.NextPageInfo, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil && args.Get(2) == nil {
		return nil, nil, nil
	}
	if args.Get(0) != nil && args.Get(1) != nil {
		return args.Get(0).([]domain.UserResponseDTO), args.Get(1).(*domain.NextPageInfo), nil
	}
	return nil, nil, args.Get(2).(lib.APIError)
}

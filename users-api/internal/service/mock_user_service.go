package service

import (
	"context"
	"errors"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) NewUser(_ context.Context, request domain.NewUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}

	if args.Get(0) != nil {
		userResponse, ok := args.Get(0).(*domain.UserResponseDTO)
		if !ok {
			return nil, lib.NewInternalServerError("failed type assertion", errors.New("internal error"))
		}

		return userResponse, nil
	}

	apiError, ok := args.Get(1).(lib.APIError)
	if !ok {
		return nil, lib.NewInternalServerError("failed type assertion", errors.New("internal error"))
	}

	return nil, apiError
}

func (m *MockUserService) UpdateUser(_ context.Context, request domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}

	if args.Get(0) != nil {
		userResponse, ok := args.Get(0).(*domain.UserResponseDTO)
		if !ok {
			return nil, lib.NewInternalServerError("failed type assertion", errors.New("internal error"))
		}

		return userResponse, nil
	}

	apiError, ok := args.Get(1).(lib.APIError)
	if !ok {
		return nil, lib.NewInternalServerError("failed type assertion ", errors.New("internal error"))
	}

	return nil, apiError
}

func (m *MockUserService) GetAllUsers(_ context.Context, request domain.FindAllUsersOptionsDTO) ([]domain.UserResponseDTO, *domain.NextPageInfo, lib.APIError) {
	args := m.Called(request)
	if args.Get(0) == nil && args.Get(1) == nil && args.Get(2) == nil {
		return nil, nil, nil
	}

	if args.Get(0) != nil && args.Get(1) != nil {
		userResponseSlice, ok1 := args.Get(0).([]domain.UserResponseDTO)
		nextPageInfo, ok2 := args.Get(1).(*domain.NextPageInfo)

		if !ok1 || !ok2 {
			return nil, nil, lib.NewInternalServerError("failed type assertion", errors.New("internal error"))
		}

		return userResponseSlice, nextPageInfo, nil
	}

	apiError, ok := args.Get(2).(lib.APIError)
	if !ok {
		return nil, nil, lib.NewInternalServerError("failed type assertion", errors.New("internal error"))
	}

	return nil, nil, apiError
}

package service

import (
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"github.com/ashtishad/ecommerce/users-api/pkg/hashpassword"
	"strings"
)

type UserService interface {
	NewUser(request domain.NewUserRequestDTO) (*domain.UserResponseDTO, error)
	UpdateUser(request domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, error)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func NewUserService(repository domain.UserRepository) *DefaultUserService {
	return &DefaultUserService{repository}
}

// NewUser first generate a salt, hashedPassword, then creates a domain model from request dto,
// then Calls the repository to save(create/update) the new user, get the user model if everything okay, otherwise returns error
// Finally returns UserResponseDTO.
func (service *DefaultUserService) NewUser(request domain.NewUserRequestDTO) (*domain.UserResponseDTO, error) {
	if err := validateCreateUserInput(request); err != nil {
		return nil, err
	}

	salt, err := hashpassword.GenerateSalt()
	if err != nil {
		return nil, err
	}

	hashedPassword := hashpassword.HashPassword(request.Password, salt)

	user := domain.User{
		Email:        strings.ToLower(request.Email),
		PasswordHash: hashedPassword,
		FullName:     request.FullName,
		Phone:        request.Phone,
		SignUpOption: request.SignUpOption,
		Status:       constants.UserStatusActive,
		Timezone:     strings.ToLower(request.Timezone),
	}

	createdUser, err := service.repo.Create(user, salt)
	if err != nil {
		return nil, err
	}

	userResponseDTO := &domain.UserResponseDTO{
		UserUUID:     createdUser.UserUUID,
		Email:        createdUser.Email,
		FullName:     createdUser.FullName,
		Phone:        createdUser.Phone,
		SignUpOption: createdUser.SignUpOption,
		Status:       createdUser.Status,
		Timezone:     createdUser.Timezone,
		CreatedAt:    createdUser.CreatedAt,
		UpdatedAt:    createdUser.UpdatedAt,
	}

	return userResponseDTO, nil
}

func (service *DefaultUserService) UpdateUser(request domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, error) {
	if err := validateUpdateUserInput(request); err != nil {
		return nil, err
	}

	user := domain.User{
		UserUUID: request.UserUUID,
		Email:    strings.ToLower(request.Email),
		FullName: request.FullName,
		Phone:    request.Phone,
		Status:   constants.UserStatusActive,
		Timezone: strings.ToLower(request.Timezone),
	}

	updatedUser, err := service.repo.Update(user)
	if err != nil {
		return nil, err
	}

	userResponseDTO := &domain.UserResponseDTO{
		UserUUID:     updatedUser.UserUUID,
		Email:        updatedUser.Email,
		FullName:     updatedUser.FullName,
		Phone:        updatedUser.Phone,
		SignUpOption: updatedUser.SignUpOption,
		Status:       updatedUser.Status,
		Timezone:     updatedUser.Timezone,
		CreatedAt:    updatedUser.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
	}

	return userResponseDTO, nil
}

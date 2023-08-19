package service

import (
	"bitbucket.org/ashtishad/as_ti/domain"
)

type UserService interface {
	NewUser(request domain.NewUserRequestDTO) (*domain.UserResponseDTO, error)
	ExistingUser(request domain.ExistingUserRequestDTO) (*domain.UserResponseDTO, error)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func NewUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repository}
}

// NewUser first converts Convert request DTO to a user domain model
// then Calls the repository to save the new user, get the user model if everything okay, otherwise returns error
// Finally converts to UserResponseDTO.
func (service DefaultUserService) NewUser(request domain.NewUserRequestDTO) (*domain.UserResponseDTO, error) {
	user := domain.User{
		Email:        request.Email,
		PasswordHash: request.Password,
		FullName:     request.FullName,
		Phone:        request.Phone,
		SignUpOption: request.SignUpOption,
		Status:       "active",
	}

	createdUser, err := service.repo.Save(user)
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
		CreatedAt:    createdUser.CreatedAt,
		UpdatedAt:    createdUser.UpdatedAt,
	}

	return userResponseDTO, nil
}

// ExistingUser first converts Convert request DTO to a user domain model (just email and password is filled)
// then Calls the repository to save the new user, get the user model if everything okay, otherwise returns error
// Finally converts to UserResponseDTO.
func (service DefaultUserService) ExistingUser(request domain.ExistingUserRequestDTO) (*domain.UserResponseDTO, error) {
	userRequest := domain.User{
		Email:        request.Email,
		PasswordHash: request.Password,
	}

	existingUser, err := service.repo.FindExisting(userRequest.Email)
	if err != nil {
		return nil, err
	}

	userResponseDTO := &domain.UserResponseDTO{
		UserUUID:     existingUser.UserUUID,
		Email:        existingUser.Email,
		FullName:     existingUser.FullName,
		Phone:        existingUser.Phone,
		SignUpOption: existingUser.SignUpOption,
		Status:       existingUser.Status,
		CreatedAt:    existingUser.CreatedAt,
		UpdatedAt:    existingUser.UpdatedAt,
	}

	return userResponseDTO, nil
}

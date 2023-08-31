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
	GetAllUsers(request domain.FindAllUsersOptions) (*[]domain.UserResponseDTO, *domain.NextPageInfo, error)
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

	userResponseDTO := createdUser.ToUserResponseDTO()

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

	userResponseDTO := updatedUser.ToUserResponseDTO()

	return userResponseDTO, nil
}

func (service *DefaultUserService) GetAllUsers(request domain.FindAllUsersOptions) (*[]domain.UserResponseDTO, *domain.NextPageInfo, error) {
	opts := domain.FindAllUsersOptions{
		FromID:       request.FromID,
		PageSize:     request.PageSize,
		Status:       request.Status,
		SignUpOption: request.SignUpOption,
		Timezone:     request.Timezone,
	}

	users, nextPageInfo, err := service.repo.FindAll(opts)
	if err != nil {
		return nil, nil, err
	}

	var userDTOs []domain.UserResponseDTO

	for _, u := range *users {
		userResponseDTO := u.ToUserResponseDTO()
		userDTOs = append(userDTOs, *userResponseDTO)
	}

	return &userDTOs, nextPageInfo, nil
}

package service

import (
	"context"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"github.com/ashtishad/ecommerce/users-api/pkg/hashpassword"
	"strings"
)

type UserService interface {
	NewUser(ctx context.Context, req domain.NewUserRequestDTO) (*domain.UserResponseDTO, lib.APIError)
	UpdateUser(ctx context.Context, req domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, lib.APIError)
	GetAllUsers(ctx context.Context, req domain.FindAllUsersOptionsDTO) ([]domain.UserResponseDTO, *domain.NextPageInfo, lib.APIError)
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
func (service *DefaultUserService) NewUser(ctx context.Context, req domain.NewUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	if err := validateCreateUserInput(req); err != nil {
		return nil, lib.NewBadRequestError("invalid create user input").Wrap(err)
	}

	salt, err := hashpassword.GenerateSalt()
	if err != nil {
		return nil, lib.NewInternalServerError("unable to generate hash", err)
	}

	hashedPassword := hashpassword.HashPassword(req.Password, salt)

	user := domain.User{
		Email:        strings.ToLower(req.Email),
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		SignUpOption: req.SignUpOption,
		Status:       constants.UserStatusActive,
		Timezone:     strings.ToLower(req.Timezone),
	}

	createdUser, apiErr := service.repo.Create(ctx, user, salt)
	if apiErr != nil {
		return nil, apiErr
	}

	userResponseDTO := createdUser.ToUserResponseDTO()

	return userResponseDTO, nil
}

func (service *DefaultUserService) UpdateUser(ctx context.Context, req domain.UpdateUserRequestDTO) (*domain.UserResponseDTO, lib.APIError) {
	if err := validateUpdateUserInput(req); err != nil {
		return nil, lib.NewBadRequestError("invalid update user input").Wrap(err)
	}

	user := domain.User{
		UserUUID: req.UserUUID,
		Email:    strings.ToLower(req.Email),
		FullName: req.FullName,
		Phone:    req.Phone,
		Status:   constants.UserStatusActive,
		Timezone: strings.ToLower(req.Timezone),
	}

	updatedUser, err := service.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	userResponseDTO := updatedUser.ToUserResponseDTO()

	return userResponseDTO, nil
}

func (service *DefaultUserService) GetAllUsers(ctx context.Context, request domain.FindAllUsersOptionsDTO) ([]domain.UserResponseDTO, *domain.NextPageInfo, lib.APIError) {
	opts, err := validateFindAllUsersOpts(request)
	if err != nil {
		return nil, nil, lib.NewBadRequestError("invalid query params").Wrap(err)
	}

	users, nextPageInfo, apiErr := service.repo.FindAll(ctx, *opts)
	if apiErr != nil {
		return nil, nil, apiErr
	}

	var userDTOs []domain.UserResponseDTO

	for _, u := range users {
		userResponseDTO := u.ToUserResponseDTO()
		userDTOs = append(userDTOs, *userResponseDTO)
	}

	return userDTOs, nextPageInfo, nil
}

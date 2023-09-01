package service

import (
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"regexp"
	"strconv"
)

// validateCreateUserInput validates the input for creating a new user.
//   - Email: Must consist of alphanumeric characters, dots, underscores, percent signs, plus signs,
//     and dashes before the @ symbol.
//     After the @ symbol, there must be a top-level domain of at least
//     two alphabetical characters.
//   - FullName: Can only consist of alphabetical characters (both uppercase and lowercase) and spaces.
//   - Phone: Must consist only of digits and must be between 10 and 15 characters in length.
//   - SignUpOption: Must be either 'general' or 'google'.
func validateCreateUserInput(input domain.NewUserRequestDTO) error {
	if matched := regexp.MustCompile(constants.EmailRegex).MatchString(input.Email); !matched {
		return fmt.Errorf("invalid email, you entered %s", input.Email)
	}

	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if matched := regexp.MustCompile(constants.FullNameRegex).MatchString(input.FullName); !matched {
		return fmt.Errorf("full name can only contain letters and spaces, you entered : %s", input.FullName)
	}

	if matched := regexp.MustCompile(constants.PhoneRegex).MatchString(input.Phone); !matched {
		return fmt.Errorf("phone must contain 10 to 15 digits, you entered: %s", input.Phone)
	}

	if matched := regexp.MustCompile(constants.SignUpOptionRegex).MatchString(input.SignUpOption); !matched {
		return fmt.Errorf("signUpOption must be 'general' or 'google', you entered: %s", input.SignUpOption)
	}

	if matched := regexp.MustCompile(constants.TimezoneRegex).MatchString(input.Timezone); !matched {
		return fmt.Errorf("timezone will be in 'UTC' or 'asia/dhaka' format, you entered: %s", input.Timezone)
	}

	return nil
}

// validateUpdateUserInput validates the input for creating a new user.
//   - Email: Must consist of alphanumeric characters, dots, underscores, percent signs, plus signs,
//     and dashes before the @ symbol.
//     After the @ symbol, there must be a top-level domain of at least
//     two alphabetical characters.
//   - FullName: Can only consist of alphabetical characters (both uppercase and lowercase) and spaces.
//   - Phone: Must consist only of digits and must be between 10 and 15 characters in length.
//   - User_id: Must be an uuid.
func validateUpdateUserInput(input domain.UpdateUserRequestDTO) error {
	if matched := regexp.MustCompile(constants.UUIDRegex).MatchString(input.UserUUID); !matched {
		return fmt.Errorf("invalid uuid, you entered %s", input.UserUUID)
	}

	if matched := regexp.MustCompile(constants.EmailRegex).MatchString(input.Email); !matched {
		return fmt.Errorf("invalid email, you entered %s", input.Email)
	}

	if matched := regexp.MustCompile(constants.FullNameRegex).MatchString(input.FullName); !matched {
		return fmt.Errorf("full name can only contain letters and spaces, you entered : %s", input.FullName)
	}

	if matched := regexp.MustCompile(constants.PhoneRegex).MatchString(input.Phone); !matched {
		return fmt.Errorf("phone must contain 10 to 15 digits, you entered: %s", input.Phone)
	}

	if matched := regexp.MustCompile(constants.TimezoneRegex).MatchString(input.Timezone); !matched {
		return fmt.Errorf("timezone will be in 'UTC' or 'asia/dhaka' format, you entered: %s", input.Timezone)
	}

	return nil
}

// validateFindAllUsersOpts validates query params, and returns domain.FindAllUsersOptions
// sets default pageSize and status if provided empty
func validateFindAllUsersOpts(input domain.FindAllUsersOptionsDTO) (*domain.FindAllUsersOptions, error) {
	opts := &domain.FindAllUsersOptions{
		PageSize: constants.DefaultPageSize,
		Status:   constants.UserStatusActive,
	}

	if input.Timezone != "" {
		if matched := regexp.MustCompile(constants.TimezoneRegex).MatchString(input.Timezone); !matched {
			return nil, fmt.Errorf("timezone must be in 'UTC' or 'Asia/Dhaka' format, you entered: %s", input.Timezone)
		}
		opts.Timezone = input.Timezone
	}

	if input.Status != "" {
		if matched := regexp.MustCompile(constants.StatusRegex).MatchString(input.Status); !matched {
			return nil, fmt.Errorf("user status must be 'active', 'inactive', or 'deleted', you entered: %s", input.Status)
		}
	}

	if input.SignUpOption != "" {
		if matched := regexp.MustCompile(constants.SignUpOptionRegex).MatchString(input.SignUpOption); !matched {
			return nil, fmt.Errorf("signUpOption must be 'general' or 'google', you entered: %s", input.SignUpOption)
		}
		opts.SignUpOption = input.SignUpOption
	}

	if input.FromIDStr != "" {
		fromID, err := strconv.Atoi(input.FromIDStr)
		if err != nil || fromID < 0 {
			return nil, fmt.Errorf("invalid FromID: must be a non-negative decimal number, you entered: %s", input.FromIDStr)
		}
		opts.FromID = fromID
	}

	if input.PageSizeStr != "" {
		pageSize, err := strconv.Atoi(input.PageSizeStr)
		if err != nil || pageSize < 20 || pageSize > 100 {
			return nil, fmt.Errorf("invalid PageSize: must be between 20 and 100, you entered: %s", input.PageSizeStr)
		}
		opts.PageSize = pageSize
	}

	return opts, nil
}

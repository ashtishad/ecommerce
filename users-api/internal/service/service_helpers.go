package service

import (
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"regexp"
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

	if input.SignUpOption != constants.SignupOptGeneral && input.SignUpOption != constants.SignUpOptGoogle {
		return fmt.Errorf("sign up option must be 'general' or 'google': %s", input.SignUpOption)
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

	return nil
}

package service

import (
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/users-api/internal/domain"
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
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(emailRegex).MatchString(input.Email); !matched {
		return fmt.Errorf("invalid email, you entered %s", input.Email)
	}

	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	fullNameRegex := `^[a-zA-Z\s]+$`
	if matched := regexp.MustCompile(fullNameRegex).MatchString(input.FullName); !matched {
		return fmt.Errorf("full name can only contain letters and spaces, you entered : %s", input.FullName)
	}

	phoneRegex := `^\d{10,15}$`
	if matched := regexp.MustCompile(phoneRegex).MatchString(input.Phone); !matched {
		return fmt.Errorf("phone must contain 10 to 15 digits, you entered: %s", input.Phone)
	}

	if input.SignUpOption != "general" && input.SignUpOption != "google" {
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
	userUUIDRegex := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	if matched := regexp.MustCompile(userUUIDRegex).MatchString(input.UserUUID); !matched {
		return fmt.Errorf("invalid uuid, you entered %s", input.UserUUID)
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(emailRegex).MatchString(input.Email); !matched {
		return fmt.Errorf("invalid email, you entered %s", input.Email)
	}

	fullNameRegex := `^[a-zA-Z\s]+$`
	if matched := regexp.MustCompile(fullNameRegex).MatchString(input.FullName); !matched {
		return fmt.Errorf("full name can only contain letters and spaces, you entered : %s", input.FullName)
	}

	phoneRegex := `^\d{10,15}$`
	if matched := regexp.MustCompile(phoneRegex).MatchString(input.Phone); !matched {
		return fmt.Errorf("phone must contain 10 to 15 digits, you entered: %s", input.Phone)
	}

	return nil
}

// validateExistingUserInput validates the input for creating a new user.
//   - Email: Must consist of alphanumeric characters, dots, underscores, percent signs, plus signs,
//     and dashes before the @ symbol.
//   - Password: Must word must be eight characters long or more
func validateExistingUserInput(input domain.ExistingUserRequestDTO) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(emailRegex).MatchString(input.Email); !matched {
		return fmt.Errorf("invalid email, you entered %s", input.Email)
	}

	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

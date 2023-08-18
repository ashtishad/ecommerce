package app

import (
	"bitbucket.org/ashtishad/as_ti/domain"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

// sanityCheck checks that all required environment variables are set.
// if any of the required variables is not defined, it prints a log message.
func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			log.Println(fmt.Sprintf("environment variable %s not defined. Terminating application...", k))
		}
	}
}

// writeResponse writes api endpoint response data and correct http status code in response.
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

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

package domain

import "time"

// UserResponseDTO has necessary fields for user response
// excludes sensitive data, such as, UserID(actual database column id) and password
type UserResponseDTO struct {
	UserUUID     string    `json:"user_uuid"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	SignUpOption string    `json:"sign_up_option"` // Enum 'general', 'google'
	Status       string    `json:"status"`         // Enum 'active', 'inactive', 'deleted'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUserRequestDTO only has required fields for creating a user
// and, excluded auto-generated fields
type NewUserRequestDTO struct {
	Email        string `json:"email" `
	Password     string `json:"password" ` // It should be handled securely and hashed before storing
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	SignUpOption string `json:"sign_up_option"` // Enum 'general', 'google', can have a default value
}

// ExistingUserRequestDTO only has email for existing user
type ExistingUserRequestDTO struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}
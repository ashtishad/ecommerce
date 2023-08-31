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
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUserRequestDTO only has required fields for creating a user
// and, excluded auto-generated fields
type NewUserRequestDTO struct {
	Email        string `json:"email" `
	Password     string `json:"password" `
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	SignUpOption string `json:"sign_up_option"` // Enum 'general', 'google', can have a default value
	Timezone     string `json:"timezone"`
}

// UpdateUserRequestDTO only has required fields for updating a user
// and, excluded password, sign_up_option field
type UpdateUserRequestDTO struct {
	UserUUID string `json:"user_id"` // path param
	Email    string `json:"email" `  // while updating, email should be unique
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Timezone string `json:"timezone"`
}

// ToUserResponseDTO converts a User to a UserResponseDTO
func (u *User) ToUserResponseDTO() *UserResponseDTO {
	return &UserResponseDTO{
		UserUUID:     u.UserUUID,
		Email:        u.Email,
		FullName:     u.FullName,
		Phone:        u.Phone,
		SignUpOption: u.SignUpOption,
		Status:       u.Status,
		Timezone:     u.Timezone,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

package domain

import "time"

type User struct {
	UserID       int       `json:"user_id"`
	UserUUID     string    `json:"user_uuid"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"` // It's important to handle this field securely
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	SignUpOption string    `json:"sign_up_option"` // Enum 'general', 'google'
	Status       string    `json:"status"`         // Enum 'active', 'inactive', 'deleted'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

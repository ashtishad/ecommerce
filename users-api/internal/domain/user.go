package domain

import "time"

type User struct {
	UserID       int       `json:"user_id"`
	UserUUID     string    `json:"user_uuid"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	SignUpOption string    `json:"sign_up_option"` // Enum 'general', 'google'
	Status       string    `json:"status"`         // Enum 'active', 'inactive', 'deleted'
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

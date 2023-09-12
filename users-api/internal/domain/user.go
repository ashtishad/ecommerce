package domain

import "time"

type User struct {
	UserID       int       `json:"userId"`
	UserUUID     string    `json:"userUuid"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	FullName     string    `json:"fullName"`
	Phone        string    `json:"phone"`
	SignUpOption string    `json:"signUpOption"` // Enum 'general', 'google'
	Status       string    `json:"status"`       // Enum 'active', 'inactive', 'deleted'
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
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

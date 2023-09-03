package domain

const (
	ErrCreatingUser          = "error creating user"
	ErrUpdateUser            = "error updating user"
	ErrInsertUserIDSalt      = "error inserting user id and salt"
	ErrUserNotFound          = "user not found"
	ErrUsersNotFound         = "users not found"
	ErrScanningData          = "error scanning user data"
	ErrCheckUserByEmail      = "unexpected error on checking user exists"
	ErrUserAlreadyExistEmail = "user already exists with this email"
)

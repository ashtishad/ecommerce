package domain

import "github.com/ashtishad/ecommerce/lib"

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	Create(user User, salt string) (*User, lib.APIError) // create(unique:email)
	Update(user User) (*User, lib.APIError)              // update(unique:uuid)
	FindAll(opts FindAllUsersOptions) ([]User, *NextPageInfo, lib.APIError)
	findUserByID(userID int) (*User, lib.APIError)
	findUserByUUID(userUUID string) (*User, lib.APIError) // helper needed in update
	checkUserExistWithEmail(email string) lib.APIError
}

// FindAllUsersOptions is filters for FindAll Users
type FindAllUsersOptions struct {
	FromID       int
	PageSize     int
	Status       string
	SignUpOption string
	Timezone     string
}

type NextPageInfo struct {
	HasNextPage bool
	StartCursor int
	EndCursor   int
	TotalCount  int
}

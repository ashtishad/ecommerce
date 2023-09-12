package domain

import (
	"context"

	"github.com/ashtishad/ecommerce/lib"
)

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	Create(ctx context.Context, user User, salt string) (*User, lib.APIError) // create(unique:email)
	Update(ctx context.Context, user User) (*User, lib.APIError)              // update(unique:uuid)
	FindAll(ctx context.Context, opts FindAllUsersOptions) ([]User, *NextPageInfo, lib.APIError)

	findUserByQuery(ctx context.Context, query string, arg interface{}) (*User, lib.APIError)
	findUserByID(ctx context.Context, userID int) (*User, lib.APIError)
	findUserByUUID(ctx context.Context, userUUID string) (*User, lib.APIError) // helper needed in update
	checkUserExistWithEmail(ctx context.Context, email string) lib.APIError
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

package domain

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	FindExisting(email string, password string) (User, error) // for signup(general and google)
	Create(user User, salt string) (User, error)              // create(unique:email)
	Update(user User) (User, error)                           // update(unique:uuid)
	findUserByID(userID int) (User, error)
	findUserByUUID(userUUID string) (User, error) // helper needed in update
}

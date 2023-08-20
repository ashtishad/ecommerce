package domain

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	FindExisting(email string, password string) (User, error) // for signup(general and google)
	Save(user User, salt string) (User, error)                // create(unique:email) or edit(unique:uuid,email)
	isUserExist(email string) (bool, error)
	findUserByID(userID int) (User, error)
}

package domain

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	Save(user User, salt string) (User, error)
	FindExisting(string, string) (User, error)
	isUserExist(email string) (bool, error)
}

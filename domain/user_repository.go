package domain

// UserRepository is the secondary port of this architecture
// It will connect to the Database adapter or mock/stub adapter
type UserRepository interface {
	Save(user User) (User, error)
}

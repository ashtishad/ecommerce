package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepositoryDB struct {
	db *sql.DB
}

func NewUserRepositoryDB(dbClient *sql.DB) UserRepositoryDB {
	return UserRepositoryDB{dbClient}
}

// Save is responsible for creating user from fields provided in domain.NewUserRequestDTO
// return sql.ErrNoRows or internal server error if some error occurs in database side.
// To ensure data integrity, it refetch user information with the help of findUserByID method.
func (d UserRepositoryDB) Save(user User) (User, error) {
	const sqlInsertUser = `
	INSERT INTO users ( email, password_hash, full_name, phone, sign_up_option ) 
	VALUES ( ?, ?, ?, ?, ?)`

	result, err := d.db.ExecContext(context.Background(), sqlInsertUser, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption)
	if err != nil {
		return User{}, fmt.Errorf("error inserting user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil || id == 0 {
		return User{}, fmt.Errorf("error getting last inserted user ID: %v", err)
	}

	user, err = d.findUserByID(int(id))
	return user, err
}

// findUserByID takes userId and returns a single user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) findUserByID(userID int) (User, error) {
	query := `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at FROM users WHERE user_id = ?`
	row := d.db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("user not found with user_id: %d", userID)
		}
		return User{}, fmt.Errorf("error scanning user data: %v", err)
	}

	return user, nil
}

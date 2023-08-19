package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type UserRepositoryDB struct {
	db *sql.DB
	l  *log.Logger
}

func NewUserRepositoryDB(dbClient *sql.DB, l *log.Logger) UserRepositoryDB {
	return UserRepositoryDB{dbClient, l}
}

// Save is responsible for creating user(if not exist) from fields provided in domain.NewUserRequestDTO
// if user already exists in the database, then updates it,
// return sql.ErrNoRows or internal server error if some error occurs in database side.
// To ensure data integrity, it refetch user information with the help of findUserByID method.
func (d UserRepositoryDB) Save(user User) (User, error) {
	existingUser, err := d.FindExisting(user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return User{}, fmt.Errorf("error checking existing user: %v", err)
	}

	// using email, as it's a unique constraint and index also
	if existingUser.Email == user.Email {
		// User exists, so update
		_, err = d.db.ExecContext(context.Background(), sqlUpdateUser, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption, user.Email)
		if err != nil {
			return User{}, fmt.Errorf("error updating user: %v", err)
		}
		return d.findUserByID(existingUser.UserID)
	}

	// User doesn't exist, so insert
	result, err := d.db.ExecContext(context.Background(), sqlInsertUser, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption)
	if err != nil {
		return User{}, fmt.Errorf("error inserting user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil || id == 0 {
		return User{}, fmt.Errorf("error getting last inserted user ID: %v", err)
	}

	return d.findUserByID(int(id))
}

// FindExisting takes user email and returns existing user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) FindExisting(email string) (User, error) {
	row := d.db.QueryRow(findExistingUserByEmail, email)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, err
		}
		return User{}, fmt.Errorf("error scanning user data: %v", err)
	}

	return user, nil
}

// findUserByID takes userId and returns a single user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) findUserByID(userID int) (User, error) {
	row := d.db.QueryRow(sqlFindUserByID, userID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//d.l.Println(err.Error())
			return User{}, err
		}
		return User{}, fmt.Errorf("error scanning user data: %v", err)
	}

	return user, nil
}

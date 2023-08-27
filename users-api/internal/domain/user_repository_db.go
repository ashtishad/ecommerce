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

func (d UserRepositoryDB) Create(user User, salt string) (User, error) {
	exists, err := d.isUserExist(user.Email)
	if err != nil {
		d.l.Printf("unexpected error on checking user exists: %s", err.Error())
		return User{}, err
	}
	if exists {
		return User{}, fmt.Errorf("user already exists with this email: %s", user.Email)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return User{}, err
	}

	// Why? check commit #054e1b6d4f6dcb9d988a89f83fb39fc9b50eabe4
	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				d.l.Printf("failed to rollback in create user: %s", rollBackErr.Error())
			}
		}
	}()

	var userID int
	err = tx.QueryRowContext(context.Background(), sqlInsertUserWithReturnID, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption, user.Timezone).Scan(&userID)
	if err != nil || userID == 0 {
		return User{}, fmt.Errorf("error creating user: %w", err)
	}

	_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt)
	if err != nil {
		return User{}, err
	}

	if err = tx.Commit(); err != nil {
		return User{}, err
	}

	return d.findUserByID(userID)
}

// Update is responsible for updating a user from fields provided in domain.UpdateUserRequestDTO
// return internal server error if some error occurs in database side.
func (d UserRepositoryDB) Update(user User) (User, error) {
	existingUser, err := d.findUserByUUID(user.UserUUID)
	if err != nil {
		return User{}, err
	}

	// if a user wants to update the email, check user already exists with
	// the updated email or not
	if existingUser.Email != user.Email {
		exists, err := d.isUserExist(user.Email)
		if err != nil {
			return User{}, err
		}
		if exists {
			return User{}, fmt.Errorf("user already exist with email : %s", existingUser.Email)
		}
	}

	_, err = d.db.Exec(sqlUpdateUser, user.Email, existingUser.PasswordHash, user.FullName, user.Phone, existingUser.SignUpOption, existingUser.UserID)
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %v", err)
	}

	return d.findUserByID(existingUser.UserID)
}

// findUserByID takes userId and returns a single user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) findUserByID(userID int) (User, error) {
	row := d.db.QueryRow(sqlFindUserByID, userID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// d.l.Println(err.Error())
			return User{}, err
		}
		return User{}, fmt.Errorf("error scanning user data: %v", err)
	}

	return user, nil
}

// findUserByUUID takes userUUID and returns a single user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) findUserByUUID(userUUID string) (User, error) {
	row := d.db.QueryRow(sqlFindUserByUUID, userUUID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// d.l.Println(err.Error())
			return User{}, err
		}
		return User{}, fmt.Errorf("error scanning user data by uuid: %v", err)
	}

	return user, nil
}

// isUserExist just for quick checking user exists or not
func (d UserRepositoryDB) isUserExist(email string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(sqlIsUserExists, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}

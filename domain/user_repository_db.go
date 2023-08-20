package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/pkg/hashpassword"
	"log"
)

type UserRepositoryDB struct {
	db *sql.DB
	l  *log.Logger
}

func NewUserRepositoryDB(dbClient *sql.DB, l *log.Logger) UserRepositoryDB {
	return UserRepositoryDB{dbClient, l}
}

// Save is responsible for creating user(if not exist) and salt from fields provided in domain.NewUserRequestDTO
// if user already exists in the database, then updates user and salt.
// return sql.ErrNoRows or internal server error if some error occurs in database side.
// To ensure data integrity, it refetch user information with the help of findUserByID method.
func (d UserRepositoryDB) Save(user User, salt string) (User, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return User{}, err
	}

	exists, err := d.isUserExist(user.Email)
	if err != nil {
		return User{}, err
	}

	var userID int

	if exists {
		// user exists, so update
		_, err = tx.ExecContext(context.Background(), sqlUpdateUser, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption, user.Email)
		if err != nil {
			tx.Rollback() // Rollback transaction on error
			return User{}, fmt.Errorf("error updating user: %v", err)
		}

		// retrieve the user_id for the given email
		err = tx.QueryRow(sqlFindUserIDFromEmail, user.Email).Scan(&userID)
		if err != nil {
			tx.Rollback() // Rollback transaction on error
			return User{}, err
		}

		// update the salt in the user_salts table with the corresponding user ID
		_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt, salt)
		if err != nil {
			tx.Rollback()
			return User{}, err
		}
	} else {
		// user doesn't exist, so insert
		result, err := tx.ExecContext(context.Background(), sqlInsertUser, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption)
		if err != nil {
			tx.Rollback()
			return User{}, fmt.Errorf("error inserting user: %v", err)
		}

		id, err := result.LastInsertId()
		if err != nil || id == 0 {
			tx.Rollback()
			return User{}, fmt.Errorf("error getting last inserted user ID: %v", err)
		}
		userID = int(id)

		// insert the salt into the user_salts table with the corresponding user ID
		_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt, salt)
		if err != nil {
			tx.Rollback()
			return User{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return User{}, err
	}

	return d.findUserByID(userID)
}

// FindExisting takes user email and password from user,
// then gets existing salt, matches with new one,
// then finally returns user from email and hashed password.
func (d UserRepositoryDB) FindExisting(email string, password string) (User, error) {
	var saltHex string
	err := d.db.QueryRow(sqlFindSaltByMail, email).Scan(&saltHex)
	if err != nil {
		return User{}, err // handle error, e.g., user not found
	}

	hashedPassword := hashpassword.HashPassword(password, saltHex)

	row := d.db.QueryRow(sqlFindUserByMailAndPass, email, hashedPassword)

	var user User
	err = row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not found or incorrect password") // user not found or password mismatch
		}
		return User{}, fmt.Errorf("error scanning user data: %v", err) // other error
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

// isUserExist just for quick checking user exists or not
func (d UserRepositoryDB) isUserExist(email string) (bool, error) {
	var exists int
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?) as user_exists"
	err := d.db.QueryRow(query, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists != 0, nil
}

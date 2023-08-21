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

// Create is responsible for creating user and salt from fields provided in domain.NewUserRequestDTO
// return internal server error if some error occurs in database side.
func (d UserRepositoryDB) Create(user User, salt string) (User, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return User{}, err
	}

	result, err := tx.ExecContext(context.Background(), sqlInsertUser, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption)
	if err != nil {
		tx.Rollback()
		return User{}, fmt.Errorf("error creating user")
	}

	id, err := result.LastInsertId()
	if err != nil || id == 0 {
		tx.Rollback()
		return User{}, fmt.Errorf("error getting last inserted user ID: %v", err)
	}
	userID := int(id)

	_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt)
	if err != nil {
		tx.Rollback()
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

	_, err = d.db.Exec(sqlUpdateUser, user.Email, existingUser.PasswordHash, user.FullName, user.Phone, existingUser.SignUpOption, existingUser.UserID)
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %v", err)
	}

	return d.findUserByID(existingUser.UserID)
}

// FindExisting takes user email and original password from user,
// then checks if user found with this email, if found,
// then retrieves salt from email, then hash the input password,
// then compares hashed password with new one, if matches
// then finally return user.
func (d UserRepositoryDB) FindExisting(email string, password string) (User, error) {
	var user User
	row := d.db.QueryRow(sqlFindUserByEmail, email)
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("user with this email not found")
		}
		d.l.Printf("unexpected error happened while ")
		return User{}, err
	}

	var saltHex string
	err = d.db.QueryRow(sqlFindSaltByMail, email).Scan(&saltHex)
	if err != nil {
		d.l.Printf("user with email: %s not found in users_salt", email)
		return User{}, fmt.Errorf("user not found")
	}

	hashedPassword := hashpassword.HashPassword(password, saltHex)
	if user.PasswordHash != hashedPassword {
		d.l.Printf("incorrect password_hash for email %s", email)
		return User{}, fmt.Errorf("incorrect password, please enter the correct password")
	}

	return user, err
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

// findUserByUUID takes userUUID and returns a single user's record
// returns error if internal server error happened.
func (d UserRepositoryDB) findUserByUUID(userUUID string) (User, error) {
	row := d.db.QueryRow(sqlFindUserByUUID, userUUID)

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

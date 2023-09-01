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

func NewUserRepositoryDB(dbClient *sql.DB, l *log.Logger) *UserRepositoryDB {
	return &UserRepositoryDB{dbClient, l}
}

func (d *UserRepositoryDB) Create(user User, salt string) (*User, error) {
	exists, err := d.isUserExist(user.Email)
	if err != nil {
		d.l.Printf("unexpected error on checking user exists: %s", err.Error())
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user already exists with this email: %s", user.Email)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return d.findUserByID(userID)
}

// Update is responsible for updating a user from fields provided in domain.UpdateUserRequestDTO
// return internal server error if some error occurs in database side.
func (d *UserRepositoryDB) Update(user User) (*User, error) {
	existingUser, err := d.findUserByUUID(user.UserUUID)
	if err != nil {
		return nil, err
	}

	// if a user wants to update the email, check user already exists with
	// the updated email or not
	if existingUser.Email != user.Email {
		exists, err := d.isUserExist(user.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("user already exist with email : %s", user.Email)
		}
	}

	_, err = d.db.Exec(sqlUpdateUser, user.Email, existingUser.PasswordHash, user.FullName, user.Phone, existingUser.SignUpOption, existingUser.UserID)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return d.findUserByID(existingUser.UserID)
}

// findUserByID takes userId and returns a single user's record
// returns error if internal server error happened.
func (d *UserRepositoryDB) findUserByID(userID int) (*User, error) {
	row := d.db.QueryRow(sqlFindUserByID, userID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// d.l.Println(err.Error())
			return nil, err
		}
		return nil, fmt.Errorf("error scanning user data: %v", err)
	}

	return &user, nil
}

// findUserByUUID takes userUUID and returns a single user's record
// returns error if internal server error happened.
func (d *UserRepositoryDB) findUserByUUID(userUUID string) (*User, error) {
	row := d.db.QueryRow(sqlFindUserByUUID, userUUID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// d.l.Println(err.Error())
			return nil, err
		}
		return nil, fmt.Errorf("error scanning user data by uuid: %v", err)
	}

	return &user, nil
}

// isUserExist just for quick checking user exists or not
func (d *UserRepositoryDB) isUserExist(email string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(sqlIsUserExists, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}

// FindAll retrieves all users from the database with optional filters
func (d *UserRepositoryDB) FindAll(opts FindAllUsersOptions) ([]User, *NextPageInfo, error) {
	var users []User
	var nextPageInfo NextPageInfo

	baseQuery := "SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, timezone, created_at, updated_at FROM users WHERE user_id > $1 "
	countQuery := "SELECT COUNT(*) FROM users WHERE user_id > $1 "

	args := []interface{}{opts.FromID}
	argCount := 2

	if opts.Status != "" {
		baseQuery += fmt.Sprintf("AND status = $%d ", argCount)
		countQuery += fmt.Sprintf("AND status = $%d ", argCount)
		args = append(args, opts.Status)
		argCount++
	}

	if opts.SignUpOption != "" {
		baseQuery += fmt.Sprintf("AND sign_up_option = $%d ", argCount)
		countQuery += fmt.Sprintf("AND sign_up_option = $%d ", argCount)
		args = append(args, opts.SignUpOption)
		argCount++
	}

	if opts.Timezone != "" {
		baseQuery += fmt.Sprintf("AND timezone = $%d ", argCount)
		countQuery += fmt.Sprintf("AND timezone = $%d ", argCount)
		args = append(args, opts.Timezone)
		argCount++
	}

	baseQuery += fmt.Sprintf("LIMIT $%d", argCount)

	args = append(args, opts.PageSize)

	rows, err := d.db.Query(baseQuery, args...)
	if err != nil {
		return nil, nil, err
	}
	defer func(rows *sql.Rows) {
		rowsClsErr := rows.Close()
		if rowsClsErr != nil {
			d.l.Printf("error closing rows in find all users: %s", rowsClsErr.Error())
			return
		}
	}(rows)

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.UserID,
			&user.UserUUID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.Phone,
			&user.SignUpOption,
			&user.Status,
			&user.Timezone,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			d.l.Printf("error scanning rows in find all users %s", err.Error())
			return nil, nil, errors.New("error scanning rows in find all users")
		}
		users = append(users, user)
	}

	userCount := len(users)

	if len(users) == 0 {
		return nil, nil, errors.New("no users found")
	}

	if len(users) < opts.PageSize {
		return users, &NextPageInfo{
			HasNextPage: false,
			StartCursor: users[0].UserID,
			EndCursor:   users[userCount-1].UserID,
			TotalCount:  userCount,
		}, nil
	}

	var totalCount int
	if err = d.db.QueryRow(countQuery, args[:argCount-1]...).Scan(&totalCount); err != nil {
		d.l.Printf("error scanning rows in find all users %s", err.Error())
		return nil, nil, errors.New("error calculating total rows in find all users")
	}

	nextPageInfo.HasNextPage = totalCount > opts.PageSize
	nextPageInfo.StartCursor = users[0].UserID
	nextPageInfo.EndCursor = users[userCount-1].UserID
	nextPageInfo.TotalCount = totalCount

	return users, &nextPageInfo, nil
}

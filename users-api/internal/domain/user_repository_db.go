package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/lib"
	"log/slog"
)

type UserRepositoryDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewUserRepositoryDB(dbClient *sql.DB, l *slog.Logger) *UserRepositoryDB {
	return &UserRepositoryDB{dbClient, l}
}

func (d *UserRepositoryDB) Create(user User, salt string) (*User, lib.APIError) {
	if apiErr := d.checkUserExistWithEmail(user.Email); apiErr != nil {
		return nil, apiErr
	}

	tx, err := d.db.Begin()
	if err != nil {
		return nil, lib.NewInternalServerError("unexpected error on tx begin", err)
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				d.l.Error("failed to rollback in create user", "err", rollBackErr.Error())
				return
			}
		}
	}()

	var userID int
	err = tx.QueryRowContext(context.Background(), sqlInsertUserWithReturnID, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption, user.Timezone).Scan(&userID)
	if err != nil || userID == 0 {
		return nil, lib.NewInternalServerError("error creating user", err)
	}

	_, err = tx.Exec(sqlInsertUserIDSalt, userID, salt)
	if err != nil {
		return nil, lib.NewInternalServerError("error inserting user id and salt", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, lib.NewInternalServerError("error committing db transaction", err)
	}

	return d.findUserByID(userID)
}

// Update is responsible for updating a user from fields provided in domain.UpdateUserRequestDTO
// return internal server error if some error occurs in database side.
func (d *UserRepositoryDB) Update(user User) (*User, lib.APIError) {
	existingUser, apiErr := d.findUserByUUID(user.UserUUID)
	if apiErr != nil {
		return nil, apiErr
	}

	// if a user wants to update the email, check user already exists with
	// the updated email or not
	if existingUser.Email != user.Email {
		apiErr := d.checkUserExistWithEmail(user.Email)
		if apiErr != nil {
			return nil, apiErr
		}
	}

	_, err := d.db.Exec(sqlUpdateUser, user.Email, existingUser.PasswordHash, user.FullName, user.Phone, existingUser.SignUpOption, existingUser.UserID)
	if err != nil {
		return nil, lib.NewInternalServerError("error updating user", err)
	}

	return d.findUserByID(existingUser.UserID)
}

// findUserByID takes userId and returns a single user's record
// returns error if internal server error happened.
func (d *UserRepositoryDB) findUserByID(userID int) (*User, lib.APIError) {
	row := d.db.QueryRow(sqlFindUserByID, userID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.NewNotFoundError("user not found by id")
		}
		return nil, lib.NewInternalServerError("error scanning user data by id", err)
	}

	return &user, nil
}

// findUserByUUID takes userUUID and returns a single user's record
// returns error if internal server error happened.
func (d *UserRepositoryDB) findUserByUUID(userUUID string) (*User, lib.APIError) {
	row := d.db.QueryRow(sqlFindUserByUUID, userUUID)

	var user User
	err := row.Scan(&user.UserID, &user.UserUUID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.SignUpOption, &user.Status, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.NewNotFoundError("user not found by uuid")
		}
		return nil, lib.NewInternalServerError("error scanning user data by uuid", err)
	}

	return &user, nil
}

// checkUserExistWithEmail checks if user exist with this email or not,
// returns error if exists is true or internal server error happens,
// if exists is false then doesn't return any error.
func (d *UserRepositoryDB) checkUserExistWithEmail(email string) lib.APIError {
	var exists bool
	err := d.db.QueryRow(sqlCheckUserExistsWithEmail, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return lib.NewInternalServerError("unexpected error on checking user exists", err)
	}
	if exists == true {
		return lib.NewBadRequestError("user already exists with this email")
	}
	return nil
}

// FindAll retrieves all users from the database with optional filters
func (d *UserRepositoryDB) FindAll(opts FindAllUsersOptions) ([]User, *NextPageInfo, lib.APIError) {
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
		return nil, nil, lib.NewInternalServerError("error retrieving rows", err)
	}
	defer func(rows *sql.Rows) {
		rowsClsErr := rows.Close()
		if rowsClsErr != nil {
			d.l.Error("error closing rows in find all users", "err", rowsClsErr.Error())
			return
		}
	}(rows)

	for rows.Next() {
		var user User
		if err = rows.Scan(
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
			return nil, nil, lib.NewInternalServerError("error scanning rows", err)
		}
		users = append(users, user)
	}

	userCount := len(users)

	if userCount == 0 {
		return nil, nil, lib.NewNotFoundError("no users found")
	}

	if userCount < opts.PageSize {
		return users, &NextPageInfo{
			HasNextPage: false,
			StartCursor: users[0].UserID,
			EndCursor:   users[userCount-1].UserID,
			TotalCount:  userCount,
		}, nil
	}

	var totalCount int
	if err = d.db.QueryRow(countQuery, args[:argCount-1]...).Scan(&totalCount); err != nil {
		return nil, nil, lib.NewInternalServerError("error calculating total rows in find all users", err)
	}

	nextPageInfo.HasNextPage = totalCount > opts.PageSize
	nextPageInfo.StartCursor = users[0].UserID
	nextPageInfo.EndCursor = users[userCount-1].UserID
	nextPageInfo.TotalCount = totalCount

	return users, &nextPageInfo, nil
}

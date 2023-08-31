package domain

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// helper functions
func mockUserObj() User {
	return User{
		UserID:       1,
		UserUUID:     "some-uuid",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		Phone:        "1234567890",
		SignUpOption: "general",
		Status:       "active",
		Timezone:     "UTC",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func mockUserRows(user User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"user_id", "user_uuid", "email", "password_hash", "full_name", "phone", "sign_up_option", "status", "timezone", "created_at", "updated_at"}).
		AddRow(user.UserID, user.UserUUID, user.Email, user.PasswordHash, user.FullName, user.Phone, user.SignUpOption, user.Status, user.Timezone, user.CreatedAt, user.UpdatedAt)
}

// Utility functions to handle Regex and TrimSpace
func expectQuery(mock sqlmock.Sqlmock, query string) *sqlmock.ExpectedQuery {
	return mock.ExpectQuery(regexp.QuoteMeta(strings.TrimSpace(query)))
}

func expectExec(mock sqlmock.Sqlmock, query string) *sqlmock.ExpectedExec {
	return mock.ExpectExec(regexp.QuoteMeta(strings.TrimSpace(query)))
}

// TestIsUserExist s checking that the isUserExist function correctly constructs and runs an SQL query to check
// whether a user with a given email exists. It's testing that the function runs without errors, and that it returns
// the correct result for the given input. By using a mock database connection, the test can run without needing
// access to an actual database, and it can make sure the function is interacting with the database as expected.
func TestIsUserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepositoryDB(db, nil)

	email := "test@example.com"

	expectQuery(mock, sqlIsUserExists).WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.isUserExist(email)

	require.NoError(t, err)
	require.True(t, exists)
}

func TestFindUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, nil)

	// Test case 1: user exists
	userID := 1
	mockUser := User{
		UserID:       userID,
		UserUUID:     "some-uuid",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		Phone:        "1234567890",
		SignUpOption: constants.SignupOptGeneral,
		Status:       constants.UserStatusActive,
		Timezone:     "asia/dhaka",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{"user_id", "user_uuid", "email", "password_hash", "full_name", "phone", "sign_up_option", "status", "timezone", "created_at", "updated_at"}).
		AddRow(mockUser.UserID, mockUser.UserUUID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.Phone, mockUser.SignUpOption, mockUser.Status, mockUser.Timezone, mockUser.CreatedAt, mockUser.UpdatedAt)

	// Test case 1: user exists
	expectQuery(mock, sqlFindUserByID).WithArgs(userID).WillReturnRows(rows)

	user, err := repo.findUserByID(userID)

	require.NoError(t, err)
	require.Equal(t, mockUser, *user)

	// Test case 2: user does not exist
	userID = 999
	expectQuery(mock, sqlFindUserByID).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	user, err = repo.findUserByID(userID)

	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows))
	require.Nil(t, user)

	// Test case 3: internal error occurs
	userID = 500
	expectQuery(mock, sqlFindUserByID).WithArgs(userID).WillReturnError(errors.New("internal error"))

	user, err = repo.findUserByID(userID)
	expectedError := errors.New("error scanning user data: internal error")

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
	require.Nil(t, user)
}

func TestFindUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, nil)

	// Test case 1: user exists
	UserUUID := "83f9ecdf-a838-4892-982f-ad34d42b1480"
	mockUser := User{
		UserID:       1,
		UserUUID:     UserUUID,
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		Phone:        "1234567890",
		SignUpOption: constants.SignupOptGeneral,
		Status:       constants.UserStatusActive,
		Timezone:     "asia/dhaka",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{"user_id", "user_uuid", "email", "password_hash", "full_name", "phone", "sign_up_option", "status", "timezone", "created_at", "updated_at"}).
		AddRow(mockUser.UserID, mockUser.UserUUID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.Phone, mockUser.SignUpOption, mockUser.Status, mockUser.Timezone, mockUser.CreatedAt, mockUser.UpdatedAt)

	// Test case 1: user exists
	expectQuery(mock, sqlFindUserByUUID).WithArgs(UserUUID).WithArgs(UserUUID).WillReturnRows(rows)

	user, err := repo.findUserByUUID(UserUUID)

	require.NoError(t, err)
	require.Equal(t, mockUser, *user)

	// Test case 2: user does not exist
	UserUUID = "7b96a2fb-3fdf-43a6-b09a-a82169286fdf"
	expectQuery(mock, sqlFindUserByUUID).WithArgs(UserUUID).WillReturnError(sql.ErrNoRows)

	user, err = repo.findUserByUUID(UserUUID)

	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows))
	require.Nil(t, user)

	// Test case 3: internal error occurs
	UserUUID = "da7ccd97-686e-444c-93c6-6bef23e6a401"
	expectQuery(mock, sqlFindUserByUUID).WithArgs(UserUUID).WillReturnError(errors.New("internal error"))

	user, err = repo.findUserByUUID(UserUUID)
	expectedError := errors.New("error scanning user data by uuid: internal error")

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
	require.Nil(t, user)
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, log.New(os.Stdout, "test: ", log.LstdFlags))

	// test case 1: user created successfully
	t.Run("User created successfully", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"
		rows := mockUserRows(mockUser)

		expectQuery(mock, sqlIsUserExists).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		expectQuery(mock, sqlInsertUserWithReturnID).
			WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.Phone, mockUser.SignUpOption, mockUser.Timezone).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
		expectExec(mock, sqlInsertUserIDSalt).WithArgs(mockUser.UserID, salt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		expectQuery(mock, sqlFindUserByID).WithArgs(1).WillReturnRows(rows)

		createdUser, err := repo.Create(mockUser, salt)
		require.NoError(t, err)
		require.NotNil(t, createdUser)

		require.Equal(t, mockUser.Email, createdUser.Email)
		require.Equal(t, mockUser.FullName, createdUser.FullName)
		require.Equal(t, mockUser.Phone, createdUser.Phone)
		require.Equal(t, mockUser.SignUpOption, createdUser.SignUpOption)
		require.Equal(t, mockUser.Status, createdUser.Status)
		require.Equal(t, mockUser.Timezone, createdUser.Timezone)
	})

	// test case 2: user already exists
	t.Run("User already exists", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"

		expectQuery(mock, sqlIsUserExists).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		_, err = repo.Create(mockUser, salt)
		require.Error(t, err)
	})

	// test case 3: database error during user creation
	t.Run("Database error during user creation", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"

		expectQuery(mock, sqlIsUserExists).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		expectQuery(mock, sqlInsertUserWithReturnID).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		_, err = repo.Create(mockUser, salt)
		require.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, log.New(os.Stdout, "test: ", log.LstdFlags))

	// test case 1: user updated successfully
	t.Run("User updated successfully", func(t *testing.T) {
		existingUser := mockUserObj()
		updateUser := existingUser
		updateUser.Email = "new@example.com"
		rows := mockUserRows(updateUser)

		expectQuery(mock, sqlFindUserByUUID).WithArgs(existingUser.UserUUID).WillReturnRows(mockUserRows(existingUser))
		expectQuery(mock, sqlIsUserExists).WithArgs(updateUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		expectExec(mock, sqlUpdateUser).
			WithArgs(updateUser.Email, existingUser.PasswordHash, updateUser.FullName, updateUser.Phone, existingUser.SignUpOption, existingUser.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectQuery(mock, sqlFindUserByID).WithArgs(existingUser.UserID).WillReturnRows(rows)

		updatedUser, err := repo.Update(updateUser)
		require.NoError(t, err)
		require.Equal(t, updateUser.Email, updatedUser.Email)
	})

	// test case 2: user does not exist
	t.Run("User does not exist", func(t *testing.T) {
		nonExistingUser := mockUserObj()

		expectQuery(mock, sqlFindUserByUUID).WithArgs(nonExistingUser.UserUUID).WillReturnError(sql.ErrNoRows)

		_, err := repo.Update(nonExistingUser)
		require.Error(t, err)
	})

	// test case 3: email already exists
	t.Run("Email already exists", func(t *testing.T) {
		existingUser := mockUserObj()
		updateUser := existingUser
		updateUser.Email = "new@example.com"

		expectQuery(mock, sqlFindUserByUUID).WithArgs(existingUser.UserUUID).WillReturnRows(mockUserRows(existingUser))
		expectQuery(mock, sqlIsUserExists).WithArgs(updateUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		_, err := repo.Update(updateUser)
		require.Error(t, err)
	})
}

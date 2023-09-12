package domain

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/stretchr/testify/require"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

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

// TestCheckUserExistWithEmail tests the checkUserExistWithEmail method
// It covers three scenarios:
// 1. The user already exists with the given email.
// 2. The user does not exist with the given email.
// 3. An internal server error occurs while executing the query.
func TestCheckUserExistWithEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewUserRepositoryDB(db, logger)

	t.Run("User exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
			WithArgs("existing@email.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		apiErr := repo.checkUserExistWithEmail(context.Background(), "existing@email.com")
		require.NotNil(t, apiErr)
		require.Equal(t, "user already exists with this email", apiErr.AsMessage())
	})

	t.Run("User does not exist", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
			WithArgs("new@email.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		apiErr := repo.checkUserExistWithEmail(context.Background(), "new@email.com")
		require.Nil(t, apiErr)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
			WithArgs("error@email.com").
			WillReturnError(errors.New("some internal error"))

		apiErr := repo.checkUserExistWithEmail(context.Background(), "error@email.com")
		require.NotNil(t, apiErr)
		require.Equal(t, lib.UnexpectedDatabaseErr, apiErr.AsMessage())
	})
}

// TestFindUserByID tests the findByID method of UserRepositoryDB,
// checks correct user domain struct outputs and http status code if error occurred
// It covers the following scenarios:
// 1. Successful retrieval of a user by ID.
// 2. User not found.
// 3. Internal server error.
func TestFindUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, testLogger)

	t.Run("Find user by id successful", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		user, apiErr := repo.findByID(context.Background(), 1)
		require.Nil(t, apiErr)
		require.Equal(t, mockUser, *user)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id = \\$1").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		user, apiErr := repo.findByID(context.Background(), 2)
		require.NotNil(t, apiErr)
		require.Equal(t, lib.UnexpectedDatabaseErr, apiErr.AsMessage())
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
		require.Nil(t, user)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id = \\$1").
			WithArgs(3).
			WillReturnError(errors.New("some internal error"))

		user, apiErr := repo.findByID(context.Background(), 3)
		require.NotNil(t, apiErr)
		require.Equal(t, lib.UnexpectedDatabaseErr, apiErr.AsMessage())
		require.Equal(t, http.StatusInternalServerError, apiErr.StatusCode())
		require.Nil(t, user)
	})
}

// TestFindUserByUUID tests the findByUUID method of UserRepositoryDB,
// checks correct user domain struct outputs and http status code if error occurred
// It covers the following scenarios:
// 1. Successful retrieval of a user by UUID.
// 2. User not found.
// 3. Internal server error.
func TestFindUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, testLogger)

	t.Run("Find user by uuid successful", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_uuid = \\$1").
			WithArgs(mockUser.UserUUID).
			WillReturnRows(rows)

		user, apiErr := repo.findByUUID(context.Background(), mockUser.UserUUID)
		require.Nil(t, apiErr)
		require.Equal(t, mockUser, *user)
	})

	t.Run("User not found", func(t *testing.T) {
		UserUUID := "some-uuid"
		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_uuid = \\$1").
			WithArgs(UserUUID).
			WillReturnError(sql.ErrNoRows)

		user, apiErr := repo.findByUUID(context.Background(), UserUUID)
		require.NotNil(t, apiErr)
		require.Equal(t, lib.UnexpectedDatabaseErr, apiErr.AsMessage())
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
		require.Nil(t, user)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		UserUUID := "some-uuid"
		mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id = \\$1").
			WithArgs(UserUUID).
			WillReturnError(errors.New("error scanning user data by uuid"))

		user, apiErr := repo.findByUUID(context.Background(), UserUUID)
		require.NotNil(t, apiErr)
		require.Equal(t, lib.UnexpectedDatabaseErr, apiErr.AsMessage())
		require.Equal(t, http.StatusInternalServerError, apiErr.StatusCode())
		require.Nil(t, user)
	})
}

// TestCreate performs unit tests with mocking on the Create method of UserRepositoryDB.
//
// Test Scenarios:
// - User created successfully: Validates that the function correctly creates a new user, commits the transaction, and returns the created user object.
// - User already exists: Tests that the function returns an error when attempting to create a user with an email that already exists.
// - Database error during user creation: Ensures that the function returns an error when a database error occurs, and rolls back the transaction.
//
// Each test case uses sqlmock to simulate database interactions,
// and the require package for assertions to ensure that the behavior is as expected.
func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, testLogger)

	t.Run("User created successfully", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"
		rows := mockUserRows(mockUser)

		expectQuery(mock, sqlCheckUserExistsWithEmail).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		expectQuery(mock, sqlInsertUserWithReturnID).
			WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.Phone, mockUser.SignUpOption, mockUser.Timezone).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
		expectExec(mock, sqlInsertUserIDSalt).WithArgs(1, salt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		expectQuery(mock, sqlFindUserByID).WithArgs(1).WillReturnRows(rows)

		createdUser, err := repo.Create(context.Background(), mockUser, salt)
		require.NoError(t, err)
		require.NotNil(t, createdUser)

		require.Equal(t, mockUser.Email, createdUser.Email)
		require.Equal(t, mockUser.FullName, createdUser.FullName)
		require.Equal(t, mockUser.Phone, createdUser.Phone)
		require.Equal(t, mockUser.SignUpOption, createdUser.SignUpOption)
		require.Equal(t, mockUser.Status, createdUser.Status)
		require.Equal(t, mockUser.Timezone, createdUser.Timezone)
	})

	t.Run("User already exists", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"
		expectQuery(mock, sqlCheckUserExistsWithEmail).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		_, err = repo.Create(context.Background(), mockUser, salt)
		require.Error(t, err)
	})

	t.Run("Database error during user creation", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"

		expectQuery(mock, sqlCheckUserExistsWithEmail).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		expectQuery(mock, sqlInsertUserWithReturnID).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		_, err = repo.Create(context.Background(), mockUser, salt)
		require.Error(t, err)
	})
}

// TestUpdate performs unit tests with mocking  on the Update method of UserRepositoryDB.
//
// Test Scenarios:
// - User updated successfully: Validates that the function correctly updates an existing user and returns the updated user object.
// - User does not exist: Tests that the function returns an error when attempting to update a non-existing user.
// - Email already exists: Ensures that the function returns an error when attempting to update an email to one that already exists in the database.
//
// Each test case uses sqlmock to simulate database interactions,
// and the require package for assertions to ensure that the behavior is as expected.
func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, testLogger)

	t.Run("User updated successfully", func(t *testing.T) {
		existingUser := mockUserObj()
		updateUser := existingUser
		updateUser.Email = "new@example.com"
		rows := mockUserRows(updateUser)

		mock.ExpectBegin()
		expectQuery(mock, sqlFindUserByUUID).WithArgs(existingUser.UserUUID).WillReturnRows(mockUserRows(existingUser))
		expectQuery(mock, sqlCheckUserExistsWithEmail).WithArgs(updateUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		expectExec(mock, sqlUpdateUser).
			WithArgs(updateUser.Email, existingUser.PasswordHash, updateUser.FullName, updateUser.Phone, existingUser.SignUpOption, existingUser.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		expectQuery(mock, sqlFindUserByID).WithArgs(existingUser.UserID).WillReturnRows(rows)

		updatedUser, err := repo.Update(context.Background(), updateUser)
		require.NoError(t, err)
		require.Equal(t, updateUser.Email, updatedUser.Email)
	})

	t.Run("User does not exist", func(t *testing.T) {
		nonExistingUser := mockUserObj()

		mock.ExpectBegin()
		expectQuery(mock, sqlFindUserByUUID).WithArgs(nonExistingUser.UserUUID).WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err := repo.Update(context.Background(), nonExistingUser)
		require.Error(t, err)
		require.Equal(t, 404, err.StatusCode())
	})

	t.Run("Email already exists", func(t *testing.T) {
		existingUser := mockUserObj()
		updateUser := existingUser
		updateUser.Email = "new@example.com"

		mock.ExpectBegin()
		expectQuery(mock, sqlFindUserByUUID).WithArgs(existingUser.UserUUID).WillReturnRows(mockUserRows(existingUser))
		expectQuery(mock, sqlCheckUserExistsWithEmail).WithArgs(updateUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		mock.ExpectRollback()

		_, err := repo.Update(context.Background(), updateUser)
		require.Error(t, err)
	})
}

// TestFindAll performs unit tests with mocking on the FindAll method of UserRepositoryDB.
//
// Test Scenarios:
// - Two Filters applied successfully, FromID and PageSize: Verifies that the function works as expected when only FromID and PageSize are used as filters.
// - Users filtered by FromID, PageSize and Status: Tests the function's behavior when FromID, PageSize, and Status are used as filters.
// - All Filters Applied: Checks that the function can handle multiple filters (FromID, PageSize, Status, SignUpOption, and Timezone) simultaneously and return the expected result.
// - Empty Result with Filters: Verifies that the function returns an error and empty results when the filters do not match any users.
// - No users found: Tests that the function returns an error when there are no users in the database.
// - Negative PageSize: Ensures that the function returns an error when the PageSize is negative.
//
// Each test case uses sqlmock to simulate database interactions,
// and the require package for assertions to ensure that the behavior is as expected.
func TestFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := UserRepositoryDB{db: db, l: testLogger}

	t.Run("Two Filters applied successfully, FromID and PageSize", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(0, 1).WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT(.+) FROM users WHERE").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		opts := FindAllUsersOptions{
			FromID:   0,
			PageSize: 1,
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.NoError(t, err)
		require.Len(t, users, 1)
		require.Equal(t, false, pageInfo.HasNextPage)
		require.Equal(t, 1, pageInfo.StartCursor)
		require.Equal(t, 1, pageInfo.EndCursor)
		require.Equal(t, 1, pageInfo.TotalCount)
	})

	t.Run("Users filtered by FromID, PageSize and Status", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(0, "active", 1).WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT(.+) FROM users WHERE").WithArgs(0, "active").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		opts := FindAllUsersOptions{
			FromID:   0,
			PageSize: 1,
			Status:   "active",
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.NoError(t, err)
		require.Len(t, users, 1)
		require.Equal(t, false, pageInfo.HasNextPage)
		require.Equal(t, 1, pageInfo.StartCursor)
		require.Equal(t, 1, pageInfo.EndCursor)
		require.Equal(t, 1, pageInfo.TotalCount)
	})

	t.Run("All Filters Applied", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(0, "active", "email", "UTC", 1).WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT(.+) FROM users WHERE").WithArgs(0, "active", "email", "UTC").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		opts := FindAllUsersOptions{
			FromID:       0,
			PageSize:     1,
			Status:       "active",
			SignUpOption: "email",
			Timezone:     "UTC",
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.NoError(t, err)
		require.Len(t, users, 1)
		require.Equal(t, false, pageInfo.HasNextPage)
		require.Equal(t, 1, pageInfo.StartCursor)
		require.Equal(t, 1, pageInfo.EndCursor)
		require.Equal(t, 1, pageInfo.TotalCount)
	})

	t.Run("Empty Result with Filters", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(0, "inactive", 1).WillReturnRows(sqlmock.NewRows([]string{}))

		opts := FindAllUsersOptions{
			FromID:   0,
			PageSize: 1,
			Status:   "inactive",
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.Error(t, err)
		require.Nil(t, users)
		require.Nil(t, pageInfo)
	})

	t.Run("No users found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE").WithArgs(0, 1).WillReturnRows(sqlmock.NewRows([]string{}))

		opts := FindAllUsersOptions{
			FromID:   0,
			PageSize: 1,
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.Error(t, err)
		require.Nil(t, users)
		require.Nil(t, pageInfo)
	})

	t.Run("Negative PageSize", func(t *testing.T) {
		opts := FindAllUsersOptions{
			FromID:   0,
			PageSize: -1,
		}

		users, pageInfo, err := repo.FindAll(context.Background(), opts)

		require.Error(t, err)
		require.Nil(t, users)
		require.Nil(t, pageInfo)
	})

}

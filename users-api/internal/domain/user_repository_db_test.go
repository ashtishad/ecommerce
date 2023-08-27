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

	escapedSQL := regexp.QuoteMeta(sqlIsUserExists)
	mock.ExpectQuery(escapedSQL).WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.isUserExist(email)

	require.NoError(t, err)
	require.True(t, exists)
}

func TestFindUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, nil)
	escapedSQL := regexp.QuoteMeta(sqlFindUserByID) // escaped any special character(especially $ sign)

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
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(userID).WillReturnRows(rows)

	user, err := repo.findUserByID(userID)

	require.NoError(t, err)
	require.Equal(t, mockUser, user)

	// Test case 2: user does not exist
	userID = 999
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(userID).WillReturnError(sql.ErrNoRows)

	user, err = repo.findUserByID(userID)

	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows))
	require.Equal(t, User{}, user)

	// Test case 3: internal error occurs
	userID = 500
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(userID).WillReturnError(errors.New("internal error"))

	user, err = repo.findUserByID(userID)
	expectedError := errors.New("error scanning user data: internal error")

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
	require.Equal(t, User{}, user)
}

func TestFindUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepositoryDB(db, nil)
	escapedSQL := regexp.QuoteMeta(sqlFindUserByUUID) // escaped any special character(especially $ sign)

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
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(UserUUID).WillReturnRows(rows)

	user, err := repo.findUserByUUID(UserUUID)

	require.NoError(t, err)
	require.Equal(t, mockUser, user)

	// Test case 2: user does not exist
	UserUUID = "7b96a2fb-3fdf-43a6-b09a-a82169286fdf"
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(UserUUID).WillReturnError(sql.ErrNoRows)

	user, err = repo.findUserByUUID(UserUUID)

	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows))
	require.Equal(t, User{}, user)

	// Test case 3: internal error occurs
	UserUUID = "da7ccd97-686e-444c-93c6-6bef23e6a401"
	mock.ExpectQuery("^" + escapedSQL + "$").WithArgs(UserUUID).WillReturnError(errors.New("internal error"))

	user, err = repo.findUserByUUID(UserUUID)
	expectedError := errors.New("error scanning user data by uuid: internal error")

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
	require.Equal(t, User{}, user)
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

		mock.ExpectQuery(regexp.QuoteMeta(sqlIsUserExists)).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		mock.ExpectQuery(`^`+regexp.QuoteMeta(strings.TrimSpace(sqlInsertUserWithReturnID))+`$`).
			WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.Phone, mockUser.SignUpOption, mockUser.Timezone).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
		mock.ExpectExec(regexp.QuoteMeta(sqlInsertUserIDSalt)).WithArgs(mockUser.UserID, salt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta(sqlFindUserByID)).WithArgs(1).WillReturnRows(rows)

		createdUser, err := repo.Create(mockUser, salt)
		require.NoError(t, err)
		require.Equal(t, mockUser.UserID, createdUser.UserID)
		require.Equal(t, mockUser.Email, createdUser.Email)
	})

	// test case 2: user already exists
	t.Run("User already exists", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"

		mock.ExpectQuery(regexp.QuoteMeta(sqlIsUserExists)).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		_, err = repo.Create(mockUser, salt)
		require.Error(t, err)
	})

	// test case 3: database error during user creation
	t.Run("Database error during user creation", func(t *testing.T) {
		mockUser := mockUserObj()
		salt := "some_salt"

		mock.ExpectQuery(regexp.QuoteMeta(sqlIsUserExists)).WithArgs(mockUser.Email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectBegin()
		mock.ExpectQuery(`^` + regexp.QuoteMeta(strings.TrimSpace(sqlInsertUserWithReturnID)) + `$`).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		_, err = repo.Create(mockUser, salt)
		require.Error(t, err)
	})
}

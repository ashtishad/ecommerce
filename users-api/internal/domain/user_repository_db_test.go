package domain

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	// escaping the parentheses in the expected query
	query := "SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\$1\\)"
	mock.ExpectQuery(query).WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.isUserExist(email)

	require.NoError(t, err)
	require.True(t, exists)
}

package database

import (
	"os"
	"testing"
)

func TestGetDSNString(t *testing.T) {
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "user")
	_ = os.Setenv("DB_PASSWD", "password")
	_ = os.Setenv("DB_ADDR", "host")
	_ = os.Setenv("DB_NAME", "dbname")

	expected := "postgres://user:password@host:5432/dbname?sslmode=disable&timezone=UTC"

	result := getDSNString()

	if result != expected {
		t.Errorf("getDSNString() returned %s; expected %s", result, expected)
	}

	// clean up environment variables after the test
	_ = os.Unsetenv("DB_PORT")
	_ = os.Unsetenv("DB_USER")
	_ = os.Unsetenv("DB_PASSWD")
	_ = os.Unsetenv("DB_ADDR")
	_ = os.Unsetenv("DB_NAME")
}

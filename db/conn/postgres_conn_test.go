package conn

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

// TestGetDSNString tests the GetDSNString function that constructs a DSN string for database connection.
// uses slog library for logger with using default slog handler options.
// It validates whether the function correctly composes the DSN string using environment variables.
func TestGetDSNString(t *testing.T) {
	var (
		dbUser   = "user"
		dbPasswd = "password"
		dbAddr   = "host"
		dbPort   = "5432"
		dbName   = "dbname"
	)

	_ = os.Setenv("DB_USER", dbUser)
	_ = os.Setenv("DB_PASSWD", dbPasswd)
	_ = os.Setenv("DB_ADDR", dbAddr)
	_ = os.Setenv("DB_PORT", dbPort)
	_ = os.Setenv("DB_NAME", dbName)

	expected := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC",
		dbUser, dbPasswd, dbAddr, dbPort, dbName)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	result := GetDSNString(logger).String()

	if result != expected {
		t.Errorf("getDSNString() returned %s; expected %s", result, expected)
	}

	_ = os.Unsetenv("DB_PORT")
	_ = os.Unsetenv("DB_USER")
	_ = os.Unsetenv("DB_PASSWD")
	_ = os.Unsetenv("DB_ADDR")
	_ = os.Unsetenv("DB_NAME")
}

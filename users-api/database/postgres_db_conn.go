package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

// getDSNString constructs a PostgreSQL Data Source Name (DSN) string using environment variables.
// It sets the connection parameters such as user, password, host, port, database name, timezone, and SSL mode.
// The resulting DSN string is in the format:
// "postgres://user:password@host:port/dbname?sslmode=disable&timezone=UTC"
// Returns the constructed DSN string.
func getDSNString() string {
	portInt, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWD")),
		Host:   fmt.Sprintf("%s:%d", os.Getenv("DB_ADDR"), portInt),
		Path:   os.Getenv("DB_NAME"),
	}
	q := dsn.Query()
	q.Set("timezone", "UTC")
	q.Set("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	return dsn.String()
}

// GetDbClient creates a new database connection and returns it
func GetDbClient() *sql.DB {
	dsn := getDSNString()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error pinging the database: %v", err)
	}
	log.Printf("successfully connected to database %s", dsn)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

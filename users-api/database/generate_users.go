package database

import (
	"database/sql"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"github.com/ashtishad/ecommerce/users-api/pkg/hashpassword"
	"log/slog"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// GenerateUsers populates n numbers of random users in users and users_salts tables,
// It accepts a *sql.DB instance to interact with the actual database and
// an integer n to specify how many records do we want to insert,
// The function uses a transaction to insert the users, rolling back the transaction and logging an error message if anything goes wrong.
func GenerateUsers(db *sql.DB, l *slog.Logger, n int) {
	gofakeit.Seed(0)

	tx, err := db.Begin()
	if err != nil {
		l.Warn("unexpected error on tx begin", "err", err.Error())
		return
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				l.Warn("failed to rollback", "err", rollBackErr.Error())
				return
			}
		}
	}()

	for i := 0; i < n; i++ {
		email := gofakeit.Email()
		fullName := gofakeit.Name()
		phone := gofakeit.Phone()
		password := gofakeit.Password(true, true, true, false, false, 12)

		salt, err := hashpassword.GenerateSalt()
		if err != nil {
			l.Warn("failed to generate salt", "err", err.Error())
			return
		}
		hashedPassword := hashpassword.HashPassword(password, salt)

		signUpOption := getRandomSignUpOption()
		userStatus := getRandomUserStatus()
		timezone := gofakeit.TimeZoneRegion()

		var userID int
		err = tx.QueryRow(`INSERT INTO users (email, password_hash, full_name, phone, sign_up_option, status, timezone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`,
			email, hashedPassword, fullName, phone, signUpOption, userStatus, timezone).Scan(&userID)
		if err != nil {
			l.Warn("failed to insert user", "err", err.Error())
			return
		}

		_, err = tx.Exec(`INSERT INTO user_salts (user_id, salt) VALUES ($1, $2)`, userID, salt)
		if err != nil {
			l.Warn("failed to insert salt", "err", err.Error())
			return
		}

		time.Sleep(1 * time.Millisecond) // Sleep for 1 millisecond to ensure different created_at and updated_at
	}

	if err = tx.Commit(); err != nil {
		l.Warn("failed to commit", "err", err.Error())
		return
	}

	// update users_user_id_seq
	query := "SELECT setval('users_user_id_seq', $1)"
	_, err = db.Exec(query, n)
	if err != nil {
		l.Warn("failed to update user_id_seq", "err", err.Error())
		return
	}
}

// getRandomUserStatus generates three possible statuses: "active", "inactive", or "deleted".
// The status "active" has a 80% chance of being chosen, "inactive" and "deleted" both have a 10% chance.
func getRandomUserStatus() string {
	randNumber := rand.Intn(100)

	switch {
	case randNumber < 80: // 80% chance
		return constants.UserStatusActive
	case randNumber < 90: // 10% chance
		return constants.UserStatusInactive
	default: // 10% chance
		return constants.UserStatusDeleted
	}
}

// getRandomSignUpOption generates two possible statuses: "general", "google".
// general has 65% chance of being chosen and google has 35%.
func getRandomSignUpOption() string {
	randNumber := rand.Intn(100)
	if randNumber < 65 {
		return constants.SignupOptGeneral
	}
	return constants.SignUpOptGoogle
}

package generate

import (
	"database/sql"
	"fmt"
	"github.com/ashtishad/ecommerce/users-api/pkg/constants"
	"github.com/ashtishad/ecommerce/users-api/pkg/hashpassword"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// GenerateUsers populates n numbers of random users in users and users_salts tables,
// It accepts a *sql.DB instance to interact with the actual database and
// an integer n to specify how many records do we want to insert,
// The function uses a transaction to insert the users, rolling back the transaction and logging an error message if anything goes wrong.
func GenerateUsers(db *sql.DB, n int) {
	gofakeit.Seed(0)

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tx, err := db.Begin()
	if err != nil {
		logger.Error("error starting transaction: ", err)
		return
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				logger.Error("failed to rollback: ", rollBackErr)
			}
		}
	}()

	signUpOptions := []string{"google", "general"}

	for i := 0; i < n; i++ {
		email := gofakeit.Email()
		fullName := gofakeit.Name()
		phone := gofakeit.Phone()
		password := gofakeit.Password(true, true, true, false, false, 12)

		salt, err := hashpassword.GenerateSalt()
		if err != nil {
			logger.Error("error generating salt: ", err)
			return
		}
		hashedPassword := hashpassword.HashPassword(password, salt)

		signUpOption := signUpOptions[rand.Intn(len(signUpOptions))]
		userStatus := getRandomUserStatus()
		timezone := gofakeit.TimeZoneRegion()

		var userID int
		err = tx.QueryRow(`INSERT INTO users (email, password_hash, full_name, phone, sign_up_option, status, timezone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`,
			email, hashedPassword, fullName, phone, signUpOption, userStatus, timezone).Scan(&userID)
		if err != nil {
			logger.Error("error inserting user and fetching ID: ", err)
			return
		}

		_, err = tx.Exec(`INSERT INTO user_salts (user_id, salt) VALUES ($1, $2)`, userID, salt)
		if err != nil {
			logger.Error("error inserting salt: ", err)
			return
		}

		time.Sleep(1 * time.Millisecond) // Sleep for 1 millisecond to ensure different created_at and updated_at
	}

	if err = tx.Commit(); err != nil {
		logger.Error("error committing transaction: ", err)
		return
	}

	// Update users_user_id_seq
	query := fmt.Sprintf("SELECT setval('users_user_id_seq', %d)", n)
	_, err = db.Exec(query)
	if err != nil {
		logger.Error("error updating sequence: ", err)
		return
	}
}

// getRandomUserStatus generates three possible statuses: "active", "inactive", or "deleted".
// The status "active" has a 70% chance of being chosen, "inactive" has a 15% chance, and "deleted" also has a 15% chance.
func getRandomUserStatus() string {
	randNumber := rand.Intn(100)

	switch {
	case randNumber < 70: // 70% chance
		return constants.UserStatusActive
	case randNumber < 85: // 15% chance
		return constants.UserStatusInactive
	default: // 15% chance
		return constants.UserStatusDeleted
	}
}

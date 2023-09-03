package domain

const (
	sqlInsertUserWithReturnID = `
		INSERT INTO users ( email, password_hash, full_name, phone, sign_up_option, timezone )
		VALUES ( $1, $2, $3, $4, $5, $6) RETURNING user_id`
	sqlUpdateUser = `UPDATE users SET email = $1, password_hash = $2, full_name = $3, phone = $4, sign_up_option = $5 
             WHERE user_id = $6`

	sqlFindUserByID = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, timezone, created_at, updated_at
                       FROM users WHERE user_id = $1`
	sqlFindUserByUUID = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status,timezone, created_at, updated_at 
                       FROM users WHERE user_uuid = $1`
	sqlInsertUserIDSalt         = `INSERT INTO user_salts (user_id, salt) VALUES ($1, $2)`
	sqlCheckUserExistsWithEmail = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
)

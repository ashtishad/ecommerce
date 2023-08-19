package domain

const (
	sqlInsertUser = `
	INSERT INTO users ( email, password_hash, full_name, phone, sign_up_option ) 
	VALUES ( ?, ?, ?, ?, ?)`

	sqlUpdateUser = `
		UPDATE users 
		SET password_hash = ?, full_name = ?, phone = ?, sign_up_option = ?
		WHERE email = ?`

	sqlFindUserByID = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at 
                       FROM users WHERE user_id = ?`

	findExistingUserByEmail = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at FROM users WHERE email = ?`
)

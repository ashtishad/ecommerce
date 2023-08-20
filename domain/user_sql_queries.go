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
	sqlFindUserIDFromEmail = `SELECT user_id FROM users WHERE email = ?`
	sqlInsertUserIDSalt    = `INSERT INTO user_salts (user_id, salt) VALUES (?, ?) ON DUPLICATE KEY UPDATE salt = ?`
	sqlFindSaltByMail      = `SELECT s.salt FROM user_salts s 
                          JOIN users u ON s.user_id = u.user_id 
                          WHERE u.email = ?`
	sqlFindUserByMailAndPass = `SELECT * FROM users WHERE email = ? AND password_hash = ?`
)

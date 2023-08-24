package domain

const (
	sqlInsertUser = `
	INSERT INTO users ( email, password_hash, full_name, phone, sign_up_option ) 
	VALUES ( ?, ?, ?, ?, ?)`

	sqlUpdateUser = `UPDATE users SET email = ?, password_hash = ?, full_name = ?, phone = ?, sign_up_option = ? WHERE user_id = ?`

	sqlFindUserByID = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at 
                       FROM users WHERE user_id = ?`
	sqlFindUserByUUID = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at 
                       FROM users WHERE user_uuid = ?`
	sqlInsertUserIDSalt = `INSERT INTO user_salts (user_id, salt) VALUES (?, ?)`
	sqlFindSaltByMail   = `SELECT s.salt FROM user_salts s 
                          JOIN users u ON s.user_id = u.user_id 
                          WHERE u.email = ?`
	sqlFindUserByEmail = `SELECT user_id, user_uuid, email, password_hash, full_name, phone, sign_up_option, status, created_at, updated_at 
                       FROM users WHERE email = ?`
)

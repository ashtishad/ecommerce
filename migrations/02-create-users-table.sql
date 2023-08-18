-- migrations/02-create-users-table.sql
USE users;
CREATE TABLE IF NOT EXISTS users (
                       user_id INT AUTO_INCREMENT PRIMARY KEY,
                       user_uuid CHAR(36) NOT NULL DEFAULT (UUID()),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       full_name VARCHAR(255) NOT NULL,
                       phone VARCHAR(20) NOT NULL,
                       sign_up_option ENUM('general', 'google') DEFAULT 'general',
                       status ENUM('active', 'inactive', 'deleted') DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- chosen sign_up_option as enum, it will have fixed values and will provide extensibility and better validation.
-- if we want to add another method in future, such as, Sign up with apple.

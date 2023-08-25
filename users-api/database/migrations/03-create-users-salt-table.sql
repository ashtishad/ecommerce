-- migrations/03-create-users-salt-table.sql
\c ecommerce;
CREATE TABLE IF NOT EXISTS user_salts
(
    user_id INT      NOT NULL,
    salt    CHAR(32) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id)
);

-- DROP TABLE IF EXISTS user_salts;

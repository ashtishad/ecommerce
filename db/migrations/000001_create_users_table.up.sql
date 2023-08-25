BEGIN;

CREATE TYPE user_status AS ENUM ('active', 'inactive', 'deleted');
CREATE TYPE sign_up_option AS ENUM ('general', 'google');

CREATE TABLE IF NOT EXISTS users
(
    user_id        SERIAL PRIMARY KEY,
    user_uuid      UUID         NOT NULL DEFAULT uuid_generate_v4(),
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hash  VARCHAR(255) NOT NULL,
    full_name      VARCHAR(255) NOT NULL,
    phone          VARCHAR(20)  NOT NULL,
    sign_up_option sign_up_option        DEFAULT 'general',
    status         user_status           DEFAULT 'active',
    timezone VARCHAR(255) NOT NULL DEFAULT 'UTC',
    created_at     TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP
);

COMMIT;

-- ALTER SEQUENCE users_user_id_seq RESTART WITH 1;

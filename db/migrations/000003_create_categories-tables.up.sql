BEGIN;

CREATE TYPE category_status AS ENUM ('active', 'inactive', 'deleted');

CREATE TABLE IF NOT EXISTS categories
(
    category_id   SERIAL PRIMARY KEY,
    category_uuid uuid            NOT NULL DEFAULT uuid_generate_v4(),
    name          VARCHAR(255)    NOT NULL,
    description   TEXT,
    status        category_status NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS category_relationships
(
    ancestor_id   INT REFERENCES categories (category_id),
    descendant_id INT REFERENCES categories (category_id),
    level INT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ancestor_id, descendant_id)
);

CREATE INDEX if not exists idx_category_uuid ON categories (category_uuid);
CREATE INDEX if not exists idx_category_name ON categories (name);
CREATE INDEX if not exists idx_category_status ON categories (status);

COMMIT;

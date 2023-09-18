BEGIN;

CREATE TYPE sim_type AS ENUM ('Single', 'Dual');
CREATE TYPE codec_type AS ENUM ('SBC', 'AAC', 'Aptx', 'LDAC');
CREATE TYPE strap_type AS ENUM ('Steel', 'Metal', 'Leather');
CREATE TYPE variant_type AS ENUM ('UAE', 'Official', 'USA','China', 'UK');

CREATE TABLE IF NOT EXISTS products
(
    product_id       SERIAL PRIMARY KEY,
    product_uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
    Type VARCHAR(255) REFERENCES categories (name),
    sku              VARCHAR(255) NOT NULL UNIQUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS phone_attributes
(
    product_id     INT REFERENCES products (product_id),
    price_in_cents BIGINT CHECK (price_in_cents >= 0),
    sim_type       sim_type      NOT NULL,
    storage VARCHAR(255) NOT NULL,
    variant variant_type NOT NULL default 'Official',
    stock   INT          NOT NULL default 0,
    PRIMARY KEY (product_id)
);

CREATE TABLE IF NOT EXISTS sound_equipment_attributes
(
    product_id     INT REFERENCES products (product_id),
    price_in_cents BIGINT CHECK (price_in_cents >= 0),
    codecs         codec_type[]  NOT NULL,
    color VARCHAR(255) NOT NULL,
    stock INT          NOT NULL default 0,
    PRIMARY KEY (product_id)
);

CREATE TABLE IF NOT EXISTS wearables_attributes
(
    product_id     INT REFERENCES products (product_id),
    price_in_cents BIGINT CHECK (price_in_cents >= 0),
    strap_types    strap_type[]  NOT NULL,
    color VARCHAR(255) NOT NULL,
    stock INT          NOT NULL default 0,
    PRIMARY KEY (product_id)
);

COMMIT;

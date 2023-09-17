BEGIN;

CREATE TYPE currency_type AS ENUM ('USD', 'Taka', 'INR');
CREATE TYPE sim_type AS ENUM ('Single', 'Dual');
CREATE TYPE codec_type AS ENUM ('SBC', 'AAC', 'Aptx', 'LDAC');
CREATE TYPE strap_type AS ENUM ('Steel', 'Metal', 'Leather');

CREATE TABLE IF NOT EXISTS products
(
    product_id       SERIAL PRIMARY KEY,
    product_uuid     uuid         NOT NULL DEFAULT uuid_generate_v4(),
    root_category_id INT REFERENCES categories (category_id),
    sku              VARCHAR(255) NOT NULL UNIQUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS phone_attributes
(
    product_id     INT REFERENCES products (product_id),
    price          numeric       NOT NULL,
    currency       currency_type NOT NULL,
    sim_type       sim_type      NOT NULL,
    storage        VARCHAR(255),
    variant        VARCHAR(255),
    stock          INT,
    specifications JSONB,
    PRIMARY KEY (product_id)
);

CREATE TABLE IF NOT EXISTS sound_equipment_attributes
(
    product_id     INT REFERENCES products (product_id),
    price          numeric       NOT NULL,
    currency       currency_type NOT NULL,
    codecs         codec_type[]  NOT NULL,
    color          VARCHAR(255),
    specifications JSONB,
    PRIMARY KEY (product_id)
);

CREATE TABLE IF NOT EXISTS wearables_attributes
(
    product_id     INT REFERENCES products (product_id),
    price          numeric       NOT NULL,
    currency       currency_type NOT NULL,
    strap_types    strap_type[]  NOT NULL,
    color          VARCHAR(255),
    specifications JSONB,
    PRIMARY KEY (product_id)
);

COMMIT;

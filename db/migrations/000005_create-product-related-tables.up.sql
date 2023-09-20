BEGIN;

CREATE TYPE product_status AS ENUM ('active', 'inactive', 'deleted');
CREATE TYPE color_status AS ENUM ('active', 'inactive', 'deleted');
CREATE TYPE brand_status AS ENUM ('active', 'inactive', 'deleted');


CREATE TABLE IF NOT EXISTS brands
(
    brand_id   SERIAL PRIMARY KEY,
    brand_uuid UUID         NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    name       VARCHAR(255) NOT NULL UNIQUE,
    status     brand_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS products
(
    product_id       BIGSERIAL PRIMARY KEY,
    product_uuid     UUID           NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    name             VARCHAR(50)    NOT NULL,
    basePriceInCents BIGINT         NOT NULL,
    category_id      INT            NOT NULL REFERENCES categories (category_id),
    brand_id         INT            NOT NULL REFERENCES brands (brand_id),
    status           product_status NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS colors
(
    color_id   SERIAL PRIMARY KEY,
    color_uuid UUID         NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    name       VARCHAR(255) NOT NULL,
    hex_code   VARCHAR(256) NOT NULL,
    status     color_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_colors
(
    product_color_id SERIAL PRIMARY KEY,
    product_id       BIGINT       NOT NULL REFERENCES products (product_id),
    color_id         INT          NOT NULL REFERENCES colors (color_id),
    status           color_status NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    UNIQUE (product_id, color_id)
);

COMMIT;

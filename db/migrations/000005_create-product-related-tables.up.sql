BEGIN;

-- main products table
CREATE TABLE IF NOT EXISTS products
(
    product_id   SERIAL PRIMARY KEY,
    product_uuid UUID         NOT NULL DEFAULT uuid_generate_v4(),
    Product_name varchar(256) NOT NULL,
    category_id  INT REFERENCES categories (category_id),
    product_type varchar(256) NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- generic product_attributes table
CREATE TABLE IF NOT EXISTS product_attributes
(
    attribute_id    SERIAL PRIMARY KEY,
    product_id      INT REFERENCES products (product_id),
    attribute_name  VARCHAR(255) NOT NULL,
    attribute_value TEXT         NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- index for fast look-up based on product_id and attribute_name
CREATE INDEX idx_product_attribute ON product_attributes (product_id, attribute_name);

COMMIT;

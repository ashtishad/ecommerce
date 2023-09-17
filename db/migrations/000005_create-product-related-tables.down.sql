BEGIN;

-- Drop tables in reverse order to avoid foreign key constraint violations
DROP TABLE IF EXISTS wearables_attributes;
DROP TABLE IF EXISTS sound_equipment_attributes;
DROP TABLE IF EXISTS phone_attributes;
DROP TABLE IF EXISTS products;

DROP TYPE IF EXISTS strap_type;
DROP TYPE IF EXISTS codec_type;
DROP TYPE IF EXISTS stock_status;
DROP TYPE IF EXISTS sim_type;
DROP TYPE IF EXISTS currency_type;


COMMIT;

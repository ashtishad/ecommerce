BEGIN;

DROP TABLE IF EXISTS category_relationships;
DROP TABLE IF EXISTS categories;
DROP TYPE IF EXISTS category_status;

-- drop INDEX if  exists idx_category_uuid;
-- drop INDEX if  exists idx_category_name;

COMMIT;

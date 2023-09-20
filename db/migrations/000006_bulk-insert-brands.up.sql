BEGIN;

INSERT INTO brands (brand_id, name, status)
VALUES (1, 'Apple', 'active'),
       (2, 'Samsung', 'active'),
       (3, 'Sony', 'active'),
       (4, 'Google', 'active'),
       (5, 'OnePlus', 'active'),
       (6, 'Huawei', 'active'),
       (7, 'Xiaomi', 'active'),
       (8, 'LG', 'active'),
       (9, 'Nokia', 'active'),
       (10, 'JBL', 'active'),
       (11, 'Bose', 'active'),
       (12, 'Sennheiser', 'active'),
       (13, 'Beats', 'active'),
       (14, 'Anker', 'active');

SELECT setval('brands_brand_id_seq', (SELECT max(brand_id) FROM brands));

COMMIT;

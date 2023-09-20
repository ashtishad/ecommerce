BEGIN;

DELETE
FROM brands
WHERE brand_id BETWEEN 1 AND 14;

DO
$$
    DECLARE
        max_id bigint;
    BEGIN
        SELECT COALESCE(MAX(brand_id), 0) INTO max_id FROM brands;

        IF max_id <= 14 THEN
            PERFORM setval('brands_brand_id_seq', 1, false);
        ELSE
            PERFORM setval('brands_brand_id_seq', max_id, true);
        END IF;
    END
$$;

COMMIT;

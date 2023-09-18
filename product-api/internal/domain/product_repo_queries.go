package domain

const (
	sqlInsertProductRetID = `INSERT INTO products (name, category_id, root_category_id) VALUES ($1, $2, $3) RETURNING product_id`
	sqlFindProductByID    = `SELECT product_uuid,name, category_id, root_category_id,created_at,updated_at FROM product where product_id = $1`

	sqlGetCategoryIDAndRootCategoryID = `WITH RECURSIVE find_root AS (
    SELECT ancestor_id, descendant_id, level
    FROM category_relationships
    WHERE descendant_id = (
        SELECT category_id FROM categories WHERE category_uuid = 'd8806fac-faed-4e2d-9665-2a5693ac0b34'
    )
    UNION ALL
    SELECT cr.ancestor_id, cr.descendant_id, cr.level
    FROM category_relationships cr, find_root fr
    WHERE cr.descendant_id = fr.ancestor_id
      AND cr.level <= fr.level
)
SELECT 
    (SELECT category_id FROM categories WHERE category_uuid = 'd8806fac-faed-4e2d-9665-2a5693ac0b34') AS category_id,
    (SELECT ancestor_id FROM find_root WHERE level = 0 LIMIT 1) AS root_category_id;`
)

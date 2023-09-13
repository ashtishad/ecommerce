package domain

const (
	sqlInsertCategory       = `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`
	sqlSelectCategoryName   = `SELECT name FROM categories WHERE LOWER(name) = LOWER($1)`
	sqlSelectCategoryByID   = `SELECT category_id,category_uuid,name, description,status,created_at,updated_at FROM categories where category_id= $1`
	sqlValidateUUIDGetCatID = `SELECT category_id, EXISTS(SELECT 1 FROM categories WHERE name = $2 AND status = 'active') FROM categories WHERE category_uuid = $1`

	sqlInsertWithLevelCalculation = `
	WITH parent_level AS (
    SELECT level FROM category_relationships WHERE descendant_id = $1 LIMIT 1
	)
	INSERT INTO category_relationships (ancestor_id, descendant_id, level)
	VALUES ($2, $3, COALESCE((SELECT level FROM parent_level), 0) + 1);
`

	sqlGetAllCategoriesWithHierarchy = `
WITH RECURSIVE CategoryTree AS (
    SELECT c.category_id, c.category_uuid, CAST(NULL AS uuid) AS parent_category_uuid, c.name, c.description, c.status, 1 AS level
    FROM categories c
    WHERE c.category_id NOT IN (SELECT descendant_id FROM category_relationships)
    UNION ALL
    SELECT c.category_id, c.category_uuid, ct.category_uuid AS parent_category_uuid, c.name, c.description, c.status, ct.level + 1
    FROM categories c
             INNER JOIN category_relationships cr ON c.category_id = cr.descendant_id
             INNER JOIN CategoryTree ct ON cr.ancestor_id = ct.category_id
)
SELECT category_uuid, parent_category_uuid, level, name FROM CategoryTree ORDER BY level, category_id;
`
)

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
	VALUES ($2, $3, COALESCE((SELECT level FROM parent_level), 1) + 1);
`
)

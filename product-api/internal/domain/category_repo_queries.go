package domain

const (
	sqlInsertCategory     = `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`
	sqlSelectCategoryName = `SELECT name FROM categories WHERE LOWER(name) = LOWER($1)`
	sqlSelectCategoryByID = `SELECT category_id,category_uuid,name, description,status,created_at,updated_at FROM categories where category_id= $1`
)

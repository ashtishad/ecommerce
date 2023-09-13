package lib

import "time"

const (
	TimeoutCreateUser = 100 * time.Millisecond
	TimeoutUpdateUser = 100 * time.Millisecond
	TimeoutGetUsers   = 200 * time.Millisecond

	TimeoutCreateCategory    = 100 * time.Millisecond
	TimeoutCreateSubcategory = 200 * time.Millisecond
	TimeoutGetAllCategories  = 500 * time.Millisecond
)

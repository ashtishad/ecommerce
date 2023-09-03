package lib

const (
	UnexpectedDatabaseErr     = "unexpected database error"
	ErrTxBegin                = "unexpected error on tx begin"
	ErrTxRollback             = "failed to rollback"
	ErrTxCommit               = "error committing db transaction"
	ErrRetrievingRows         = "error retrieving rows"
	ErrClosingRows            = "error closing rows"
	ErrScanningRows           = "error scanning rows"
	ErrTotalCountInPagination = "error calculating total count in pagination"
)

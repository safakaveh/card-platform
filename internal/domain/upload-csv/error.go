package uploadcsv

import "errors"

var (
	ErrMissingFile       = errors.New("CSV file is required")
	ErrInvalidCSV        = errors.New("CSV file is invalid")
	ErrNoMappedColumns   = errors.New("no frn_ or bck_ column was found")
	ErrDuplicateColumn   = errors.New("CSV contains duplicate mapped columns")
	ErrEmptyOrderName    = errors.New("order name is required")
	ErrOrderNameConflict = errors.New("an import with this order name already exists")
)

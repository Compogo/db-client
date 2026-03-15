package repository

// Page represents pagination parameters for database queries.
// It defines which subset of results to return and how many items per page.
type Page struct {
	// Number is the page number (starting from 1).
	Number uint64

	// Limit is the maximum number of items per page.
	Limit uint32
}

// NewPage creates a new Page instance with the given number and limit.
// Example:
//
//	page := NewPage(1, 20)  // first page, 20 items per page
func NewPage(number uint64, limit uint32) *Page {
	return &Page{Number: number, Limit: limit}
}

// Sort represents sorting criteria for database queries.
type Sort struct {
	// ColumnName is the name of the database column to sort by.
	ColumnName string

	// Direction is the sort direction (ASC or DESC).
	Direction SortDirection
}

// NewSort creates a new Sort instance with the given column and direction.
// Example:
//
//	sort := NewSort("created_at", DESC)
func NewSort(columnName string, direction SortDirection) *Sort {
	return &Sort{ColumnName: columnName, Direction: direction}
}

// Filter represents filtering criteria for database queries.
// It defines a condition on a specific column with a value and comparison operator.
type Filter struct {
	// ColumnName is the name of the database column to filter on.
	ColumnName string

	// Value is the value to compare against (as string, will be converted by the driver).
	Value string

	// Comparable is the comparison operator (Eq, Neq, Gt, Gte, Lt, Lte, LIKE, IN).
	Comparable Comparable
}

// NewFilter creates a new Filter instance with the given column, value, and operator.
// Example:
//
//	filter := NewFilter("age", "25", Gte)  // age >= 25
func NewFilter(columnName string, value string, comparable Comparable) *Filter {
	return &Filter{ColumnName: columnName, Value: value, Comparable: comparable}
}

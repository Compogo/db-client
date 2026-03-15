package repository

//go:generate stringer -type=Comparable

// Comparable defines comparison operators for filtering database queries.
// Each constant represents a standard SQL operator.
const (
	// Eq is the equality operator '='.
	Eq Comparable = iota

	// Neq is the inequality operator '!='.
	Neq

	// Gt is the greater than operator '>'.
	Gt

	// Gte is the greater than or equal operator '>='.
	Gte

	// Lt is the less than operator '<'.
	Lt

	// Lte is the less than or equal operator '<='.
	Lte

	// LIKE is the pattern matching operator 'LIKE'.
	// Use with '%' wildcards in the value.
	LIKE

	// IN is the membership operator 'IN(...)'.
	// The value should contain comma-separated values.
	IN
)

// Comparable represents a comparison operator for filtering.
type Comparable uint8

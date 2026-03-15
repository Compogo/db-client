package repository

//go:generate stringer -type=SortDirection -output sort_direction_string.go

// SortDirection defines the order of sorting for query results.
const (
	// ASC indicates ascending order (A to Z, smallest to largest).
	ASC SortDirection = iota

	// DESC indicates descending order (Z to A, largest to smallest).
	DESC
)

// SortDirection represents the sorting direction (ASC or DESC).
type SortDirection uint8

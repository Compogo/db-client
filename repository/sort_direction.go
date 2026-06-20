package repository

//go:generate stringer -type=SortDirection -output sort_direction_string.go

// SortDirection — направление сортировки.
const (
	ASC  SortDirection = iota // По возрастанию
	DESC                      // По убыванию
)

// SortDirection представляет направление сортировки.
type SortDirection uint8

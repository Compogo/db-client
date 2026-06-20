package repository

// Page содержит параметры пагинации.
type Page struct {
	Number uint64
	Limit  uint64
}

// NewPage создаёт новый Page.
func NewPage(number uint64, limit uint64) *Page {
	return &Page{Number: number, Limit: limit}
}

// Sort определяет сортировку по одной колонке.
type Sort struct {
	ColumnName string
	Direction  SortDirection
}

// NewSort создаёт новый Sort.
func NewSort(columnName string, direction SortDirection) *Sort {
	return &Sort{ColumnName: columnName, Direction: direction}
}

// Filter определяет условие фильтрации.
type Filter struct {
	ColumnName string
	Value      string
	Comparable Comparable
}

// NewFilter создаёт новый Filter.
func NewFilter(columnName string, value string, comparable Comparable) *Filter {
	return &Filter{ColumnName: columnName, Value: value, Comparable: comparable}
}

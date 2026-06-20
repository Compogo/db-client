package repository

//go:generate stringer -type=Comparable

// Comparable — тип сравнения для фильтров.
const (
	Eq   Comparable = iota // равно
	Neq                    // не равно
	Gt                     // больше
	Gte                    // больше или равно
	Lt                     // меньше
	Lte                    // меньше или равно
	LIKE                   // LIKE (поиск по шаблону)
	IN                     // IN (список значений)
)

// Comparable определяет операцию сравнения.
type Comparable uint8

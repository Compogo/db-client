package driver

// Driver represents a database driver identifier (e.g., "postgres", "mysql").
// It implements fmt.Stringer for consistent string representation.
type Driver string

// String returns the driver name as a string.
func (d Driver) String() string {
	return string(d)
}

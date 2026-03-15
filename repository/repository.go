package repository

import "context"

// Pager defines the interface for paginated data retrieval.
// It returns a slice of items for the requested page, applying sorting and filters.
type Pager[T any] interface {
	// Page returns a paginated list of items.
	// Parameters:
	//   - ctx: context for cancellation and timeouts
	//   - page: pagination parameters (page number and limit)
	//   - sorts: sorting criteria (may be empty)
	//   - filters: filtering criteria (may be empty)
	//
	// Returns:
	//   - []T: slice of items for the current page
	//   - error: any error encountered during the query
	Page(ctx context.Context, page *Page, sorts []*Sort, filters ...*Filter) ([]T, error)
}

// Counter defines the interface for counting items matching given filters.
type Counter interface {
	// Count returns the total number of items matching the filters.
	// This is typically used to calculate the total number of pages
	// when implementing pagination.
	//
	// Parameters:
	//   - ctx: context for cancellation and timeouts
	//   - filters: filtering criteria (may be empty)
	//
	// Returns:
	//   - int: total count of matching items
	//   - error: any error encountered during the query
	Count(ctx context.Context, filters ...*Filter) (int, error)
}

// Repository combines pagination and counting capabilities for a specific entity type.
// It provides a standard interface for data access operations that can be implemented
// by different database drivers (PostgreSQL, MySQL, etc.).
//
// Example:
//
//	type UserRepository interface {
//	    repository.Repository[User]
//	    FindByEmail(ctx context.Context, email string) (*User, error)
//	}
type Repository[T any] interface {
	Pager[T]
	Counter
}

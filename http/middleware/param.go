package middleware

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/db-client/repository"
	"github.com/Compogo/http"
	"github.com/Compogo/http/middleware/param"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
	"github.com/spf13/cast"
)

const (
	// PageFieldName is the query parameter name for page number.
	PageFieldName = "page"

	// LimitFieldName is the query parameter name for items per page.
	LimitFieldName = "limit"

	// SortFieldName is the query parameter name for sorting.
	SortFieldName = "sort"

	// FilterFieldName is the query parameter name for filtering.
	FilterFieldName = "filter"

	// CountHeaderName is the HTTP header for total item count.
	CountHeaderName = "X-Pagination-Count"

	// PageHeaderName is the HTTP header for current page number.
	PageHeaderName = "X-Pagination-Page"

	// LimitHeaderName is the HTTP header for items per page.
	LimitHeaderName = "X-Pagination-Limit"

	// SortHeaderName is the HTTP header for sorting parameters.
	SortHeaderName = "X-Pagination-Sort"

	// FilterHeaderName is the HTTP header for filter parameters.
	FilterHeaderName = "X-Pagination-Filter"

	// LimitDefault is the default items per page.
	LimitDefault = uint8(20)

	// PageDefault is the default page number.
	PageDefault = uint32(1)

	// SortDirectionAsc is the string representation of ascending sort.
	SortDirectionAsc = "ASC"

	// SortDirectionDesc is the string representation of descending sort.
	SortDirectionDesc = "DESC"

	// FilterEq is the string representation of equality operator.
	FilterEq = "EQ"

	// FilterNeq is the string representation of not equal operator.
	FilterNeq = "Neq"

	// FilterGt is the string representation of greater than operator.
	FilterGt = "Gt"

	// FilterGte is the string representation of greater than or equal operator.
	FilterGte = "Gte"

	// FilterLt is the string representation of less than operator.
	FilterLt = "Lt"

	// FilterLte is the string representation of less than or equal operator.
	FilterLte = "Lte"

	// FilterLIKE is the string representation of LIKE operator.
	FilterLIKE = "LIKE"

	// FilterIN is the string representation of IN operator.
	FilterIN = "IN"

	itemSeparator = ";"
	pairSeparator = ":"
)

var (
	// SortDirectionToQuerySort maps repository.SortDirection to query string values.
	SortDirectionToQuerySort = linker.NewLinker[repository.SortDirection, string](
		linker.NewLink(repository.ASC, SortDirectionAsc),
		linker.NewLink(repository.DESC, SortDirectionDesc),
	)

	// QuerySortToSortDirection maps query string values to repository.SortDirection.
	QuerySortToSortDirection = linker.NewLinker[string, repository.SortDirection](
		linker.NewLink(strings.ToLower(SortDirectionAsc), repository.ASC),
		linker.NewLink(strings.ToLower(SortDirectionDesc), repository.DESC),
	)

	// ComparableToComparableQuery maps repository.Comparable to query string operators.
	ComparableToComparableQuery = linker.NewLinker[repository.Comparable, string](
		linker.NewLink(repository.Eq, FilterEq),
		linker.NewLink(repository.Neq, FilterNeq),
		linker.NewLink(repository.Gt, FilterGt),
		linker.NewLink(repository.Gte, FilterGte),
		linker.NewLink(repository.Lt, FilterLt),
		linker.NewLink(repository.Lte, FilterLte),
		linker.NewLink(repository.LIKE, FilterLIKE),
		linker.NewLink(repository.IN, FilterIN),
	)

	// QueryComparableToComparable maps query string operators to repository.Comparable.
	QueryComparableToComparable = linker.NewLinker[string, repository.Comparable](
		linker.NewLink(FilterEq, repository.Eq),
		linker.NewLink(FilterNeq, repository.Neq),
		linker.NewLink(FilterGt, repository.Gt),
		linker.NewLink(FilterGte, repository.Gte),
		linker.NewLink(FilterLt, repository.Lt),
		linker.NewLink(FilterLte, repository.Lte),
		linker.NewLink(FilterLIKE, repository.LIKE),
		linker.NewLink(FilterIN, repository.IN),
	)
)

// NewSort creates middleware that parses sort parameters into []*repository.Sort.
// The sort parameter format is: "column:direction;column:direction"
// Example: "created_at:ASC;name:DESC"
//
// Parameters:
//   - logger: for logging parsing errors
//   - defaultSort: fallback sort when none is provided
//   - columns: list of allowed column names for security
//
// The parsed value is stored in request context under "sort".
func NewSort(logger logger.Logger, defaultSort *repository.Sort, columns ...string) http.Middleware {
	allowedColumns := set.NewSet[string](columns...)

	return param.NewParam(
		SortFieldName,
		logger,
		func(value any) (any, error) {
			sort, err := cast.ToStringE(value)
			if err != nil {
				return nil, err
			}

			pairs := strings.Split(sort, itemSeparator)

			sorts := make([]*repository.Sort, len(pairs))

			for i, pair := range pairs {
				values := strings.SplitN(pair, pairSeparator, 2)
				if len(values) != 2 {
					return nil, fmt.Errorf("[PaginationSort] sort invalid pair %s", pair)
				}

				if !allowedColumns.Contains(values[0]) {
					return nil, fmt.Errorf("[PaginationSort] sort invalid column %s", values[0])
				}

				sortDirection, err := QuerySortToSortDirection.Get(strings.ToLower(values[1]))
				if err != nil {
					return nil, fmt.Errorf("[PaginationSort] invalid sort direction %s", values[1])
				}

				sorts[i] = repository.NewSort(sort, sortDirection)
			}

			return sorts, nil
		},
		param.WithUriGetter(),
		param.WithHeaderGetterByName(SortHeaderName),
		param.WithDefault(defaultSort),
	)
}

// NewFilter creates middleware that parses filter parameters into []*repository.Filter.
// The filter format is: "column:value:operator;column:value:operator"
// Example: "age:25:Gt;status:active:Eq"
//
// Parameters:
//   - logger: for logging parsing errors
//   - columns: list of allowed column names for security
//
// The parsed value is stored in request context under "filter".
func NewFilter(logger logger.Logger, columns ...string) http.Middleware {
	allowedColumns := set.NewSet[string](columns...)

	return param.NewParam(
		FilterFieldName,
		logger,
		func(value any) (any, error) {
			sort, err := cast.ToStringE(value)
			if err != nil {
				return nil, err
			}

			pairs := strings.Split(sort, itemSeparator)

			filters := make([]*repository.Filter, len(pairs))

			for i, pair := range pairs {
				values := strings.SplitN(pair, pairSeparator, 3)
				if len(values) != 3 {
					return nil, fmt.Errorf("[PaginationSort] filter invalid pair %s", pair)
				}

				if !allowedColumns.Contains(values[0]) {
					return nil, fmt.Errorf("[PaginationSort] filter invalid column %s", values[0])
				}

				comparable, err := QueryComparableToComparable.Get(strings.ToLower(values[2]))
				if err != nil {
					return nil, fmt.Errorf("[PaginationSort] filter invalid comparable %s", values[2])
				}

				filters[i] = repository.NewFilter(values[0], values[1], comparable)
			}

			return filters, nil
		},
		param.WithUriGetter(),
		param.WithHeaderGetterByName(FilterHeaderName),
		param.WithDefault(""),
	)
}

// NewPage creates middleware that parses the page number parameter.
// The page parameter is validated to be >= 1.
//
// The parsed value is stored in request context under "page" as uint64.
// Default value is PageDefault (1).
func NewPage(logger logger.Logger) http.Middleware {
	return param.NewParamUint64(
		PageFieldName,
		logger,
		param.WithHeaderGetterByName(PageHeaderName),
		param.WithUriGetter(),
		param.WithDefault(PageDefault),
		param.AddValidator(param.GteValidator(1)),
	)
}

// NewPaginationLimit creates middleware that parses the items per page parameter.
// The limit parameter is validated to be >= 1.
//
// The parsed value is stored in request context under "limit" as uint32.
// Default value is LimitDefault (20).
func NewPaginationLimit(logger logger.Logger) http.Middleware {
	return param.NewParamUint32(
		LimitFieldName,
		logger,
		param.WithHeaderGetterByName(LimitHeaderName),
		param.WithUriGetter(),
		param.WithDefault(LimitDefault),
		param.AddValidator(param.GteValidator(1)),
	)
}

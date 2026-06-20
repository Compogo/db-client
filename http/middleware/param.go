package middleware

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/Compogo/db-client/repository"
	httpServer "github.com/Compogo/http_server"
	"github.com/Compogo/http_server/middleware/param"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
	"github.com/spf13/cast"
)

const (
	// PageFieldName - имя поля для номера страницы
	PageFieldName = "page"

	// LimitFieldName - имя поля для лимита
	LimitFieldName = "limit"

	// SortFieldName - имя поля для сортировки
	SortFieldName = "sort"

	// FilterFieldName - имя поля для фильтрации
	FilterFieldName = "filter"

	// CountHeaderName - заголовок с общим количеством
	CountHeaderName = "X-Pagination-Count"

	// PageHeaderName - заголовок с номером страницы
	PageHeaderName = "X-Pagination-Page"

	// LimitHeaderName - заголовок с лимитом
	LimitHeaderName = "X-Pagination-Limit"

	// SortHeaderName - заголовок с сортировкой
	SortHeaderName = "X-Pagination-Sort"

	// FilterHeaderName - заголовок с фильтрацией
	FilterHeaderName = "X-Pagination-Filter"

	// LimitDefault - лимит по умолчанию
	LimitDefault = uint64(20)

	// PageDefault - страница по умолчанию
	PageDefault = uint64(1)

	// SortDirectionAsc - сортировка по возрастанию
	SortDirectionAsc = "ASC"

	// SortDirectionDesc - сортировка по убыванию
	SortDirectionDesc = "DESC"

	// FilterEq - равно
	FilterEq = "EQ"

	// FilterNeq - не равно
	FilterNeq = "Neq"

	// FilterGt - больше
	FilterGt = "Gt"

	// FilterGte - больше или равно
	FilterGte = "Gte"

	// FilterLt - меньше
	FilterLt = "Lt"

	// FilterLte - меньше или равно
	FilterLte = "Lte"

	// FilterLIKE - LIKE
	FilterLIKE = "LIKE"

	// FilterIN - IN
	FilterIN = "IN"

	// itemSeparator разделитель элементов в строке
	itemSeparator = ";"

	// pairSeparator разделитель ключ-значение в элементе
	pairSeparator = ":"
)

var (
	// SortDirectionToQuerySort Связь между repository.SortDirection и строковым представлением.
	SortDirectionToQuerySort = linker.NewLinker[repository.SortDirection, string](
		linker.Link(repository.ASC, SortDirectionAsc),
		linker.Link(repository.DESC, SortDirectionDesc),
	)

	// QuerySortToSortDirection Связь между строковым представлением и repository.SortDirection.
	QuerySortToSortDirection = linker.NewLinker[string, repository.SortDirection](
		linker.Link(strings.ToLower(SortDirectionAsc), repository.ASC),
		linker.Link(strings.ToLower(SortDirectionDesc), repository.DESC),
		linker.KeyStringNormalizer[repository.SortDirection](),
	)

	// ComparableToComparableQuery Связь между repository.Comparable и строковым представлением.
	ComparableToComparableQuery = linker.NewLinker[repository.Comparable, string](
		linker.Link(repository.Eq, FilterEq),
		linker.Link(repository.Neq, FilterNeq),
		linker.Link(repository.Gt, FilterGt),
		linker.Link(repository.Gte, FilterGte),
		linker.Link(repository.Lt, FilterLt),
		linker.Link(repository.Lte, FilterLte),
		linker.Link(repository.LIKE, FilterLIKE),
		linker.Link(repository.IN, FilterIN),
	)

	// QueryComparableToComparable Связь между строковым представлением и repository.Comparable.
	QueryComparableToComparable = linker.NewLinker[string, repository.Comparable](
		linker.Link(FilterEq, repository.Eq),
		linker.Link(FilterNeq, repository.Neq),
		linker.Link(FilterGt, repository.Gt),
		linker.Link(FilterGte, repository.Gte),
		linker.Link(FilterLt, repository.Lt),
		linker.Link(FilterLte, repository.Lte),
		linker.Link(FilterLIKE, repository.LIKE),
		linker.Link(FilterIN, repository.IN),
		linker.KeyStringNormalizer[repository.Comparable](),
	)
)

// NewSort создаёт middleware для извлечения сортировки из запроса.
// Поддерживает формат: "column:asc;column2:desc"
func NewSort(logger compogo.Logger, defaultSort *repository.Sort, columns ...string) httpServer.Middleware {
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
					return nil, fmt.Errorf("[Pagination][Sort] sort invalid pair %s", pair)
				}

				if !allowedColumns.Contains(values[0]) {
					return nil, fmt.Errorf("[Pagination][Sort] sort invalid column %s", values[0])
				}

				sortDirection, err := QuerySortToSortDirection.Get(values[1])
				if err != nil {
					return nil, fmt.Errorf("[Pagination][Sort] invalid sort direction %s", values[1])
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

// NewFilter создаёт middleware для извлечения фильтров из запроса.
// Поддерживает формат: "column:value:EQ;column2:value2:LIKE"
func NewFilter(logger compogo.Logger, columns ...string) httpServer.Middleware {
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
					return nil, fmt.Errorf("[Pagination][Sort] filter invalid pair %s", pair)
				}

				if !allowedColumns.Contains(values[0]) {
					return nil, fmt.Errorf("[Pagination][Sort] filter invalid column %s", values[0])
				}

				comparable, err := QueryComparableToComparable.Get(values[2])
				if err != nil {
					return nil, fmt.Errorf("[Pagination][Sort] filter invalid comparable %s", values[2])
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

// NewPage создаёт middleware для извлечения номера страницы из запроса.
func NewPage(logger compogo.Logger) httpServer.Middleware {
	return param.NewParamUint64(
		PageFieldName,
		logger,
		param.WithHeaderGetterByName(PageHeaderName),
		param.WithUriGetter(),
		param.WithDefault(PageDefault),
		param.AddValidator(param.GteValidator(1)),
	)
}

// NewPaginationLimit создаёт middleware для извлечения лимита из запроса.
func NewPaginationLimit(logger compogo.Logger) httpServer.Middleware {
	return param.NewParamUint64(
		LimitFieldName,
		logger,
		param.WithHeaderGetterByName(LimitHeaderName),
		param.WithUriGetter(),
		param.WithDefault(LimitDefault),
		param.AddValidator(param.GteValidator(1)),
	)
}

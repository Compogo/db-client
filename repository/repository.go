package repository

import "context"

// Pager — интерфейс для постраничного получения данных.
type Pager[T any] interface {
	Page(context.Context, *Page, []*Sort, ...*Filter) ([]T, error)
}

// Counter — интерфейс для подсчёта количества записей.
type Counter interface {
	Count(context.Context, ...*Filter) (uint64, error)
}

// Saver — интерфейс для сохранения сущности.
// Объединяет Insert и Update, определяя операцию по наличию ID.
type Saver[T any] interface {
	Save(context.Context, T) (T, error)
}

// Inserter — интерфейс для вставки новой сущности.
type Inserter[T any] interface {
	Insert(context.Context, T) (T, error)
}

// Updater — интерфейс для обновления существующей сущности.
type Updater[T any] interface {
	Update(context.Context, T) (T, error)
}

// Deleter — интерфейс для удаления записей по фильтрам.
type Deleter interface {
	Delete(context.Context, ...*Filter) error
}

// Finder — интерфейс для поиска записей по фильтрам.
type Finder[T any] interface {
	Find(context.Context, ...*Filter) ([]T, error)
}

// BulkInserter — интерфейс для массовой вставки сущностей.
type BulkInserter[T any] interface {
	BulkInsert(context.Context, ...T) error
}

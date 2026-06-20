# Compogo DB Client

[![Go Reference](https://pkg.go.dev/badge/github.com/Compogo/db-client.svg)](https://pkg.go.dev/github.com/Compogo/db-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Плагинная система для работы с базами данных в фреймворке [Compogo](https://github.com/Compogo/compogo).

Предоставляет:

* Единый интерфейс для работы с любой БД (MySQL, PostgreSQL, SQLite)
* Плагинную систему драйверов
* Логирование всех SQL-запросов
* Лимитер для защиты от ошибок соединения
* Middleware для пагинации, сортировки и фильтрации

## Установка

```shell
go get github.com/Compogo/db-client
```

## Быстрый старт

```go
package main

import (
    "context"
    "github.com/Compogo/compogo"
    "github.com/Compogo/db-client"
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithComponents(&db_client.Component),
    )

    app.AddComponents(&compogo.Component{
        Name: "user_service",
        Init: compogo.StepFunc(func(container compogo.Container) error {
            return container.Invoke(func(client db_client.Client) error {
                rows, err := client.Query("SELECT * FROM users WHERE active = ?", true)
                if err != nil {
                    return err
                }
                defer rows.Close()
                // обработка rows
                return nil
            })
        }),
    })

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Использование

```go
import (
    "github.com/Compogo/db-client"
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithComponents(&db_client.Component),
    )
}
```

## Конфигурация через флаги

```shell
# Выбор драйвера (автоматически подставляет доступные)
--db.driver=mysql
```

## Логирование запросов

```go
// Автоматически логирует все запросы на уровне Debug
// [Database][mysql] query: SELECT * FROM users; args: [true]
```

## Лимитер ошибок соединения

```go
// Автоматически ограничивает количество ошибок соединения
// При превышении лимита переходит в режим ожидания
```

## Repository

### Интерфейсы для работы с данными

```go
// Постраничное получение
type Pager[T any] interface {
    Page(context.Context, *Page, []*Sort, ...*Filter) ([]T, error)
}

// Подсчёт количества
type Counter interface {
    Count(context.Context, ...*Filter) (uint64, error)
}

// Сохранение (Insert или Update)
type Saver[T any] interface {
    Save(context.Context, T) (T, error)
}

// Вставка
type Inserter[T any] interface {
    Insert(context.Context, T) (T, error)
}

// Обновление
type Updater[T any] interface {
    Update(context.Context, T) (T, error)
}

// Удаление
type Deleter interface {
    Delete(context.Context, ...*Filter) error
}

// Поиск
type Finder[T any] interface {
    Find(context.Context, ...*Filter) ([]T, error)
}

// Массовая вставка
type BulkInserter[T any] interface {
    BulkInsert(context.Context, ...T) error
}
```

### Middleware для пагинации

```go
import "github.com/Compogo/db-client/middleware"

// Использование в роутере
router.Use(
    middleware.NewPage(logger),      // извлекает номер страницы
    middleware.NewPaginationLimit(logger), // извлекает лимит
    middleware.NewSort(logger, defaultSort, "id", "name", "created_at"), // сортировка
    middleware.NewFilter(logger, "name", "age", "active"), // фильтрация
)

// Формат сортировки: ?sort=name:asc;created_at:desc
// Формат фильтрации: ?filter=name:john:LKE;age:18:Gt
```

## API

### Client

```go
type Client interface {
    io.Closer
    Query(string, ...interface{}) (*sql.Rows, error)
    QueryRow(string, ...interface{}) *sql.Row
    Exec(string, ...interface{}) (sql.Result, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    SQL() *sql.DB
    DriverName() string
}
```

### Config

```go
type Config struct {
    Driver string // имя драйвера БД
}
```

### Repository

```go
type Comparable uint8 // Eq, Neq, Gt, Gte, Lt, Lte, LIKE, IN

type SortDirection uint8 // ASC, DESC

type Page struct {
    Number uint64
    Limit  uint64
}

type Sort struct {
    ColumnName string
    Direction  SortDirection
}

type Filter struct {
    ColumnName string
    Value      string
    Comparable Comparable
}
```

## Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [Compogo Runner](https://github.com/Compogo/runner) — для Limiter
* [database/sql](https://pkg.go.dev/database/sql) — стандартная библиотека

## Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```

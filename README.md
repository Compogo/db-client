# Compogo DB Client 💾

**Compogo DB Client** — это гибкий и расширяемый клиент для работы с базами данных, построенный на принципах плагинной архитектуры. Позволяет подключать различные драйверы (PostgreSQL, MySQL, SQLite и др.) через единый интерфейс, добавлять функциональность через декораторы и полностью интегрируется с жизненным циклом Compogo.

## 🚀 Установка

```bash
go get github.com/Compogo/db-client
```

### 📦 Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/db-client"
    _ "github.com/Compogo/postgres" // импортируем нужный драйвер
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithOsSignalCloser(),
        db_client.Component, // базовый компонент БД
        compogo.WithComponents(
            userRepositoryComponent,
        ),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}

// Компонент, использующий БД
var userRepositoryComponent = &component.Component{
    Dependencies: component.Components{db_client.Component},
    Execute: component.StepFunc(func(c container.Container) error {
        return c.Invoke(func(db db_client.Client) {
            // db готов к работе
            rows, _ := db.Query("SELECT * FROM users")
            defer rows.Close()
            // ...
        })
    }),
}
```

### ✨ Возможности

#### 🎯 Плагинная архитектура драйверов

Любой драйвер регистрируется одной строкой и автоматически становится доступным:

```go
// В драйвере postgres
func init() {
    db_client.Registration(Postgres, NewPostgresClient)
}
```

Выбор драйвера через флаг:

```bash
./myapp --db.driver=postgres
```

#### 🔌 Единый интерфейс Client

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
    Driver() Driver
}
```

Полностью совместим со стандартным `database/sql`.

### 🎨 Декораторы (паттерн Decorator)

Добавляйте функциональность, не изменяя основной код:

#### Логирование запросов

```go
import "github.com/Compogo/db-client/logger"

db = &logger.Logger{Client: db, logger: logger}
// Все запросы теперь логируются на уровне DEBUG
```

#### Circuit Breaker (защита от сбоев)

```go
import "github.com/Compogo/db-client/connection"

db = connection.NewLimiter(db, 5, 5*time.Second)
// После 5 ошибок подряд — 5 секунд "тишины"
```
### Пример драйвера (postgres)

```go
var Component = &component.Component{
    Init: component.StepFunc(func(c container.Container) error {
        return container.Provides(
            NewConfig,
            NewClient,
        )
    }),
    // ... другие шаги
}

func init() {
    db_client.Registration(Postgres, NewClientFromContainer)
}
```

### Репозиторий с Circuit Breaker и логгированием

```go
type UserRepository struct {
    db db_client.Client
}

func NewUserRepository(db db_client.Client) *UserRepository {
    // Оборачиваем клиент
    db = &logger.Logger{Client: db, logger: logger}
    db = connection.NewLimiter(db, 3, 10*time.Second)
    
    return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*User, error) {
    row := r.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id)
    var user User
    err := row.Scan(&user.ID, &user.Name)
    return &user, err
}
```

### 🔧 Создание своего драйвера

```go
const Postgres db_client.Driver = "postgres"

type postgresClient struct {
    db *sql.DB
    // ...
}

func NewPostgresClient(container container.Container) (db_client.Client, error) {
    // достаём конфиг из контейнера
    var config *Config
    if err := container.Invoke(func(cfg *Config) { config = cfg }); err != nil {
        return nil, err
    }
    
    // создаём подключение
    db, err := sql.Open(Postgres.String(), config.DSN)
    if err != nil {
        return nil, err
    }
    
    return &postgresClient{db: db}, nil
}

func init() {
    db_client.Registration(Postgres, NewPostgresClient)
}
```

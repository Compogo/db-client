package db_client

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/types/linker"
)

// drivers — хранилище зарегистрированных драйверов БД.
// Ключ — имя драйвера, значение — функция создания клиента.
var drivers = linker.NewLinker[string, Getter](linker.KeyStringNormalizer[Getter]())

// Registration регистрирует драйвер БД.
// Должна вызываться в init() каждого пакета драйвера.
//
// Пример регистрации MySQL-драйвера:
//
//	func init() {
//	    db_client.Registration("mysql", func(container compogo.Container) (db_client.Client, error) {
//	        var config *mysql.Config
//	        container.Invoke(func(c *mysql.Config) { config = c })
//	        return mysql.NewClient(config)
//	    })
//	}
func Registration(driverName string, getter Getter) {
	drivers.Add(driverName, getter)
}

// Getter — фабричная функция для создания клиента БД.
// Принимает DI-контейнер для получения зависимостей драйвера.
type Getter func(container compogo.Container) (Client, error)

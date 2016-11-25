package mesql

import (
	"database/sql"
	"database/sql/driver"

	_ "github.com/go-sql-driver/mysql"
)

var drivers = map[string]*driver.Driver{}
var dbs = map[string]*DB{}

// 注册并命名驱动
func RegisterDriver(name string, driver driver.Driver) error {
	if drivers[name] != nil {
		return ErrRegister
	}
	if driver == nil {
		return ErrDriver
	}
	sql.Register(name, driver)
	drivers[name] = &driver
	return nil
}

// 注册数据库连接
func RegisterDB(name, driverName, dataSourceName string) error {
	var db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	dbs[name] = &DB{db: db, autocommit: true}
	dbs[name].Init()
	return nil
}

// 打开连接
func OpenDB(name string) *DB {
	var db, ok = dbs[name]
	if ok {
		return db
	}
	return nil
}

// init.go 文件进行数据库初始化操作
package medb

import (
	"database/sql"
	"database/sql/driver"
)

var (
	// 驱动映射
	drivers = map[string]*driver.Driver{}
	// 数据库连接
	dbs = map[string]*DB{}
)

// 注册并命名驱动
// 如果同名驱动已经注册，返回错误;否则将驱动加入到驱动映射中
func RegisterDriver(name string, driver driver.Driver) error {
	if drivers[name] != nil {
		panic(ErrRegister)
	}
	if driver == nil {
		return ErrDriver
	}
	sql.Register(name, driver)
	drivers[name] = &driver
	return nil
}

// 注册数据库连接
// name:给数据库连接的命名
// driverName:驱动名
// dataSourceName：数据库连接信息
func RegisterDB(name, driverName, dataSourceName string) error {
	if dbs[name] != nil {
		return ErrRegisterDB
	}
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

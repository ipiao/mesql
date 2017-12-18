package medb

import (
	"database/sql"
	"errors"
	"sync"
	"time"
)

var (
	dbs            = map[string]*DB{}
	maxOpenConnNum = 30
	maxIdleConnNum = 10
	maxLifeTime    = time.Minute * 30
	showSQL        = false
)

// 标签
const (
	MedbTag       = "db"
	MedbFieldName = "col"
)

// ShowSQL 是否打印日志
func ShowSQL(b bool) {
	showSQL = b
}

// RegisterDB 注册数据库连接
// name:给数据库连接的命名
// driverName:驱动名
// dataSourceName：数据库连接信息
func RegisterDB(name, driverName, dataSourceName string) error {
	var mu = sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	if dbs[name] != nil {
		return errors.New("连接已存在")
	}
	var db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxOpenConnNum)
	db.SetMaxIdleConns(maxIdleConnNum)
	db.SetConnMaxLifetime(maxLifeTime)

	dbs[name] = &DB{DB: db, name: name}
	return db.Ping()
}

// OpenDB 打开连接
func OpenDB(name string) *DB {
	var db, ok = dbs[name]
	if ok {
		return db
	}
	return nil
}

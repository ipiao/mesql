package medb

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ipiao/metools/logger"
)

var (
	dbs            = map[string]*DB{}
	maxOpenConnNum = 30
	maxIdleConnNum = 10
	maxLifeTime    = time.Minute * 30
)

// exported for changing supported
var (
	MedbTag           = "db"        // medb 标签字段
	MedbFieldName     = "col"       // 标签解析映射后，col对应字段名
	MedbFieldIgnore   = "_"         // "_" not "-"
	MedbFieldCp       = "cp"        // custome parse 字段自定义解析标签
	MedbFieldCpMethod = "MedbParse" // custome parse 字段自定义解析方法名
	Logger            = melogger.New("medb")
)

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

// ParseTag 解析标签
// TODO mem cached
func ParseTag(tag string) map[string]string {
	var res = make(map[string]string)
	var arr = strings.Split(tag, ",")
	for _, a := range arr {
		if strings.Contains(a, ":") {
			brr := strings.Split(a, ":")
			res[brr[0]] = brr[1]
		} else {
			res[a] = ""
		}
	}
	if len(res) == 1 {
		if _, ok := res[MedbFieldIgnore]; ok {
			res[MedbFieldName] = MedbFieldIgnore
		} else {
			for t, v := range res {
				if v == "" {
					res[MedbFieldName] = t
				}
			}
		}
	}
	return res
}

func logSQL(err error, sql string, args ...interface{}) {
	if err != nil {
		Logger.Errorj(map[string]interface{}{"sql": sql, "args": args, "error": err})
	} else {
		Logger.Debugj(map[string]interface{}{"sql": sql, "args": args})
	}
}

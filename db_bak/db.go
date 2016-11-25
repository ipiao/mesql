package mesql

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"time"
)

// 自定义DB，包含db连接信息
type DB struct {
	db         *sql.DB
	tx         *sql.Tx
	autocommit bool
}

// 初始化db
func (this *DB) Init() {
	this.db.SetMaxOpenConns(DefaultMaxOpenConns)
	this.db.SetMaxIdleConns(DefaultMaxIdleConns)
	this.db.SetConnMaxLifetime(DefaultConnMaxLifetime)
}

// DB的驱动
func (this *DB) Driver() driver.Driver {
	return this.db.Driver()
}

// Ping检查连接是否有效
func (this *DB) Ping() error {
	return this.db.Ping()
}

// Close,关闭连接
func (this *DB) Close() error {
	return this.db.Close()
}

// 设置最大连接数
func (this *DB) SetMaxOpenConns(n int) {
	this.db.SetMaxOpenConns(n)
}

// 设置最大空闲连接数
func (this *DB) SetMaxIdleConns(n int) {
	this.db.SetMaxIdleConns(n)
}

// 设置连接时间
func (this *DB) SetConnMaxLifetime(d time.Duration) {
	this.db.SetConnMaxLifetime(d)
}

// 解析sql
func (this *DB) Exec(sql string, params ...interface{}) (sql.Result, error) {
	log.Println("[tinysql]:", sql)
	if this.autocommit {
		return this.db.Exec(sql, params...)
	} else {
		return this.tx.Exec(sql, params...)
	}
}

// 查询
func (this *DB) Query(sql string, params ...interface{}) *Rows {
	log.Println("[tinysql]:", sql)
	if this.autocommit {
		var rows, err = this.db.Query(sql, params...)
		return &Rows{rows: rows, err: err}

	} else {
		var rows, err = this.tx.Query(sql, params...)
		return &Rows{rows: rows, err: err}
	}

}

// 这个方法没什么意义
func (this *DB) queryRow(sql string, params ...interface{}) *sql.Row {
	if this.autocommit {
		return this.db.QueryRow(sql, params...)
	} else {
		return this.tx.QueryRow(sql, params...)
	}
}

//
func (this *DB) Prepare(sql string) *Stmt {
	if this.autocommit {
		var stmt, err = this.db.Prepare(sql)
		return &Stmt{stmt: stmt, err: err}
	} else {
		var stmt, err = this.tx.Prepare(sql)
		return &Stmt{stmt: stmt, err: err}
	}
}

// 开启事务
func (this *DB) Begin() bool {
	var err error
	if !this.autocommit {
		this.tx, err = this.db.Begin()
	}
	if err != nil {
		return false
	}
	return true
}

// 提交事务
func (this *DB) Commit() error {
	if this.autocommit {
		return ErrCommit
	}
	var err = this.tx.Commit()
	if err == nil {
		this.autocommit = true
	}
	return err
}

// 回滚
func (this *DB) RollBack() error {
	if this.autocommit {
		return ErrRollBack
	}
	var err = this.tx.Rollback()
	if err == nil {
		this.autocommit = true
	}
	return err
}

// 使已经存在的状态生成事务的状态
func (this *DB) Stmt(stmt *Stmt) *Stmt {
	var s = this.tx.Stmt(stmt.stmt)
	return &Stmt{stmt: s}
}

//
func (this *DB) Stats() sql.DBStats {
	return this.db.Stats()
}

// 连接数
func (this *DB) OpenConnetcions() int {
	return this.db.Stats().OpenConnections
}

// NewBuilder 创建sql构造器
func (this *DB) NewBuilder() *builder {
	var b = new(builder)
	b.reset()
	b.db = this
	return b
}

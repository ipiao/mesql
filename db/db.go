package medb

import (
	"database/sql"
	"database/sql/driver"
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

// 解析sql
func (this *DB) Exec(sql string, params ...interface{}) (sql.Result, error) {
	if this.autocommit {
		return this.db.Exec(sql, params...)
	} else {
		return this.tx.Exec(sql, params...)
	}
}

// 查询
func (this *DB) Query(sql string, params ...interface{}) *Rows {
	if this.autocommit {
		var rows, err = this.db.Query(sql, params...)
		return &Rows{rows: rows, err: err}

	} else {
		var rows, err = this.tx.Query(sql, params...)
		return &Rows{rows: rows, err: err}
	}

}

// 查询最多单行
func (this *DB) QueryRow(sql string, params ...interface{}) *Row {
	if this.autocommit {
		var row = this.db.QueryRow(sql, params...)
		return &Row{row: row}
	} else {
		var row = this.tx.QueryRow(sql, params...)
		return &Row{row: row}
	}
}

// 预处理
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

// 返回连接的信息
func (this *DB) Stats() sql.DBStats {
	return this.db.Stats()
}

// 连接数
func (this *DB) OpenConnetcions() int {
	return this.db.Stats().OpenConnections
}

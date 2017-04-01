package medb

import (
	"database/sql"
	"errors"
	"sync"
)

// DB 自定义DB
type DB struct {
	*sql.DB
	*sql.Tx
	autoCommit bool
}

// MountDB 嵌入db
func (d *DB) MountDB(db *sql.DB) error {
	var mu = new(sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if d.DB != nil {
		return errors.New("db already has connection")
	}
	d.DB = db
	d.autoCommit = true
	return nil
}

// Exec 解析sql
func (d *DB) Exec(sql string, args ...interface{}) *Result {
	if !d.autoCommit {
		var res, err = d.Tx.Exec(sql, args...)
		return &Result{res, err}
	}
	var res, err = d.DB.Exec(sql, args...)
	return &Result{res, err}
}

// Query 查询
func (d *DB) Query(sql string, args ...interface{}) *Rows {
	if !d.autoCommit {
		var rows, err = d.Tx.Query(sql, args...)
		return &Rows{Rows: rows, err: err}
	}
	var rows, err = d.DB.Query(sql, args...)
	return &Rows{Rows: rows, err: err}
}

// QueryRow 查询单行
func (d *DB) QueryRow(sql string, args ...interface{}) *Row {
	if !d.autoCommit {
		var row = d.Tx.QueryRow(sql, args...)
		return &Row{row}
	}
	var row = d.DB.QueryRow(sql, args...)
	return &Row{row}
}

// Prepare 预处理
func (d *DB) Prepare(sql string) *Stmt {
	if !d.autoCommit {
		var stmt, err = d.Tx.Prepare(sql)
		return &Stmt{Stmt: stmt, err: err}
	}
	var stmt, err = d.DB.Prepare(sql)
	return &Stmt{Stmt: stmt, err: err}
}

// Begin 开启事务
func (d *DB) Begin() error {
	var err error
	if d.autoCommit {
		d.Tx, err = d.DB.Begin()
		if err != nil {
			return err
		}
		d.autoCommit = false
	}
	return nil
}

// Commit 提交事务
func (d *DB) Commit() error {
	if d.autoCommit {
		return errors.New("transaction closed")
	}
	d.autoCommit = true
	var err = d.Tx.Commit()
	return err
}

// RollBack 回滚
func (d *DB) RollBack() error {
	if d.autoCommit {
		return errors.New("transaction closed")
	}
	d.autoCommit = true
	var err = d.Tx.Rollback()
	return err
}

// Stmt 使已经存在的状态生成事务的状态
func (d *DB) Stmt(stmt *Stmt) *Stmt {
	if d.autoCommit {
		return &Stmt{err: errors.New("transaction closed")}
	}
	var s = d.Tx.Stmt(stmt.Stmt)
	return &Stmt{Stmt: s}
}

// OpenConnetcions 连接数
func (d *DB) OpenConnetcions() int {
	return d.DB.Stats().OpenConnections
}

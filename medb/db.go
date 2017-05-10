package medb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
)

// DB 自定义DB
type DB struct {
	*sql.DB
	name string
}

// MountDB 嵌入db
func (d *DB) MountDB(db *sql.DB, name string) error {
	var mu = new(sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if d.DB != nil {
		return errors.New("db already has connection")
	}
	d.DB = db
	d.name = name
	return nil
}

// Exec 解析sql
func (d *DB) Exec(sql string, args ...interface{}) *Result {
	var res, err = d.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("[medb] sql exec error:sql='%s',args=%v", sql, args)
	}
	return &Result{res, err}
}

// ExecContext 解析sql
func (d *DB) ExecContext(ctx context.Context, sql string, args ...interface{}) *Result {
	if !d.autoCommit {
		var res, err = d.Tx.ExecContext(ctx, sql, args...)
		if err != nil {
			log.Printf("[medb] tx sql exec context error:sql='%s',args=%v", sql, args)
		}
		return &Result{res, err}
	}
	var res, err = d.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		log.Printf("[medb] sql exec context error:sql='%s',args=%v", sql, args)
	}
	return &Result{res, err}
}

// Query 查询
func (d *DB) Query(sql string, args ...interface{}) *Rows {
	if !d.autoCommit {
		var rows, err = d.Tx.Query(sql, args...)
		if err != nil {
			log.Printf("[medb] tx sql query error:sql='%s',args=%v", sql, args)
		}
		return &Rows{Rows: rows, err: err}
	}
	var rows, err = d.DB.Query(sql, args...)
	if err != nil {
		log.Printf("[medb] sql query error:sql='%s',args=%v", sql, args)
	}
	return &Rows{Rows: rows, err: err}
}

// QueryContext 查询
func (d *DB) QueryContext(ctx context.Context, sql string, args ...interface{}) *Rows {
	if !d.autoCommit {
		var rows, err = d.Tx.QueryContext(ctx, sql, args...)
		if err != nil {
			log.Printf("[medb] tx sql query context error:sql='%s',args=%v", sql, args)
		}
		return &Rows{Rows: rows, err: err}
	}
	var rows, err = d.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		log.Printf("[medb] sql query context error:sql='%s',args=%v", sql, args)
	}
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

// PrepareContext 预处理
func (d *DB) PrepareContext(ctx context.Context, sql string) *Stmt {
	if !d.autoCommit {
		var stmt, err = d.Tx.PrepareContext(ctx, sql)
		return &Stmt{Stmt: stmt, err: err}
	}
	var stmt, err = d.DB.PrepareContext(ctx, sql)
	return &Stmt{Stmt: stmt, err: err}
}

// Begin 开启事务
func (d *DB) Begin() error {
	var err error
	d.Tx, err = d.DB.Begin()
	if err != nil {
		return err
	}
	d.autoCommit = false
	return nil
}

// BeginTx 开启事务
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	var err error
	d.Tx, err = d.DB.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	d.autoCommit = false
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

// Rollback 回滚
func (d *DB) Rollback() error {
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

package medb

import (
	"context"
	"database/sql"
	"errors"
)

// DB 自定义DB
type DB struct {
	*sql.DB
	name string
}

// Name 返回注册名
func (d *DB) Name() string {
	return d.name
}

// MountDB 嵌入db
func (d *DB) MountDB(db *sql.DB, name string) error {
	if d.DB != nil {
		return errors.New("db already has connection")
	}
	d.DB = db
	d.name = name
	return nil
}

// Exec 解析sql
func (d *DB) Exec(sql string, args ...interface{}) *Result {
	return d.ExecContext(context.TODO(), sql, args...)
}

// ExecContext 解析sql
func (d *DB) ExecContext(ctx context.Context, sql string, args ...interface{}) *Result {
	var res, err = d.DB.ExecContext(ctx, sql, args...)
	logSQL(err, sql, args...)
	return &Result{res, err}
}

// Query 查询
func (d *DB) Query(sql string, args ...interface{}) *Rows {
	return d.QueryContext(context.TODO(), sql, args...)
}

// QueryContext 查询
func (d *DB) QueryContext(ctx context.Context, sql string, args ...interface{}) *Rows {
	var rows, err = d.DB.QueryContext(ctx, sql, args...)
	logSQL(err, sql, args...)
	return &Rows{Rows: rows, err: err}
}

// QueryRow 查询单行
func (d *DB) QueryRow(sql string, args ...interface{}) *Row {
	var row = d.DB.QueryRow(sql, args...)
	return &Row{row}
}

// Prepare 预处理
func (d *DB) Prepare(sql string) *Stmt {
	var stmt, err = d.DB.Prepare(sql)
	return &Stmt{Stmt: stmt, err: err}
}

// PrepareContext 预处理
func (d *DB) PrepareContext(ctx context.Context, sql string) *Stmt {
	var stmt, err = d.DB.PrepareContext(ctx, sql)
	return &Stmt{Stmt: stmt, err: err}
}

// Begin 开启事务
func (d *DB) Begin() (*Tx, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// BeginTx 开启事务
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// OpenConnetcions 连接数
func (d *DB) OpenConnetcions() int {
	return d.Stats().OpenConnections
}

// Commit 提交事务
func (d *DB) Commit() error {
	return nil
}

// Rollback 回滚
func (d *DB) Rollback() error {
	return nil
}

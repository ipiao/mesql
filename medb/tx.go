package medb

import (
	"context"
	"database/sql"
	"log"
)

// Tx 事务
type Tx struct {
	*sql.Tx
}

// Exec 解析sql
func (d *Tx) Exec(sql string, args ...interface{}) *Result {
	var res, err = d.Tx.Exec(sql, args...)
	logSQL(err, sql, args...)
	return &Result{res, err}
}

// ExecContext 解析sql
func (d *Tx) ExecContext(ctx context.Context, sql string, args ...interface{}) *Result {
	var res, err = d.Tx.ExecContext(ctx, sql, args...)
	if err != nil {
		log.Printf("[medb] tx sql exec context error:sql='%s',args=%v", sql, args)
	}
	return &Result{res, err}
}

// Query 查询
func (d *Tx) Query(sql string, args ...interface{}) *Rows {
	var rows, err = d.Tx.Query(sql, args...)
	logSQL(err, sql, args...)
	return &Rows{Rows: rows, err: err}
}

// QueryContext 查询
func (d *Tx) QueryContext(ctx context.Context, sql string, args ...interface{}) *Rows {
	var rows, err = d.Tx.QueryContext(ctx, sql, args...)
	logSQL(err, sql, args...)
	return &Rows{Rows: rows, err: err}
}

// QueryRow 查询单行
func (d *Tx) QueryRow(sql string, args ...interface{}) *Row {
	var row = d.Tx.QueryRow(sql, args...)
	return &Row{row}
}

// Prepare 预处理
func (d *Tx) Prepare(sql string) *Stmt {
	var stmt, err = d.Tx.Prepare(sql)
	return &Stmt{Stmt: stmt, err: err}
}

// PrepareContext 预处理
func (d *Tx) PrepareContext(ctx context.Context, sql string) *Stmt {
	var stmt, err = d.Tx.PrepareContext(ctx, sql)
	return &Stmt{Stmt: stmt, err: err}
}

// Stmt 使已经存在的状态生成事务的状态
func (d *Tx) Stmt(stmt *Stmt) *Stmt {
	var s = d.Tx.Stmt(stmt.Stmt)
	return &Stmt{Stmt: s}
}

// Commit 提交事务
func (d *Tx) Commit() error {
	var err = d.Tx.Commit()
	return err
}

// Rollback 回滚
func (d *Tx) Rollback() error {
	var err = d.Tx.Rollback()
	return err
}

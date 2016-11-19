package medb

import (
	"database/sql"
)

// 准备状态
// 使用时注意手动关闭
type Stmt struct {
	stmt *sql.Stmt
	err  error
}

// 返回错误信息
func (this *Stmt) Error() error {
	return this.err
}

// 解析
func (this *Stmt) Exec(params ...interface{}) *Result {
	var res, err = this.stmt.Exec(params...)
	return &Result{result: res, err: err}
}

// 查询
func (this *Stmt) Query(params ...interface{}) *Rows {
	var rows, err = this.stmt.Query(params...)
	return &Rows{rows: rows, err: err}
}

// 查询单行
func (this *Stmt) QueryRow(params ...interface{}) *Row {
	var row = this.stmt.QueryRow(params...)
	return &Row{row: row}
}

// 关闭
func (this *Stmt) Close() error {
	return this.stmt.Close()
}

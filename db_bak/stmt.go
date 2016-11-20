package mesql

import (
	"database/sql"
)

// 准备状态
type Stmt struct {
	stmt *sql.Stmt
	err  error
}

// 返回错误信息
func (this *Stmt) Error() error {
	return this.err
}

// 解析
func (this *Stmt) Exec(params ...interface{}) (sql.Result, error) {
	return this.stmt.Exec(params...)
}

// 查询
func (this *Stmt) Query(params ...interface{}) *Rows {
	var rows, err = this.stmt.Query(params...)
	return &Rows{rows: rows, err: err}
}

// 查询至多单行
func (this *Stmt) QueryRow(params ...interface{}) *sql.Row {
	return this.stmt.QueryRow(params...)
}

// 关闭
func (this *Stmt) Close() error {
	return this.stmt.Close()
}

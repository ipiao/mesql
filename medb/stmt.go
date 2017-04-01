package medb

import (
	"database/sql"
)

// Stmt 准备状态
// 使用时注意手动关闭
type Stmt struct {
	*sql.Stmt
	err error
}

// Err 返回错误信息
func (s *Stmt) Err() error {
	return s.err
}

// Exec 解析
func (s *Stmt) Exec(params ...interface{}) *Result {
	var res, err = s.Stmt.Exec(params...)
	return &Result{Result: res, err: err}
}

// Query 查询
func (s *Stmt) Query(params ...interface{}) *Rows {
	var rows, err = s.Stmt.Query(params...)
	return &Rows{Rows: rows, err: err}
}

// QueryRow 查询单行
func (s *Stmt) QueryRow(params ...interface{}) *Row {
	var row = s.Stmt.QueryRow(params...)
	return &Row{row}
}

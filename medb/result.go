package medb

import (
	"database/sql"
)

// Result 执行结果
type Result struct {
	sql.Result
	err error
}

// Err 返回错误
func (r *Result) Err() error {
	return r.err
}

// SetErr 错误设置
func (r *Result) SetErr(err error) *Result {
	r.err = err
	return r
}

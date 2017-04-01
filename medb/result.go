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

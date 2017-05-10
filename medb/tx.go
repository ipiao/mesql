package medb

import (
	"database/sql"
	"log"
)

// Tx 事务
type Tx struct {
	*sql.Tx
	err error
}

// Exec 解析sql
func (d *Tx) Exec(sql string, args ...interface{}) *Result {
	var res, err = d.Exec(sql, args...)
	if err != nil {
		log.Printf("[medb] tx sql exec error:sql='%s',args=%v", sql, args)
	}
	return &Result{res, err}
}

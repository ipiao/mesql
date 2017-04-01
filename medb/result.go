package medb

import (
	"database/sql"
)

// Result 执行结果
type Result struct {
	sql.Result
	err error
}

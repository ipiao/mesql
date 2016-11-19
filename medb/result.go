package medb

import (
	"database/sql"
)

type Result struct {
	result sql.Result
	err    error
}

func (this *Result) Error() error {
	return this.err
}

func (this *Result) LastInsertIdInt() int {
	var n, _ = this.result.LastInsertId()
	return int(n)
}
func (this *Result) RowsAffectedInt() int {
	var n, _ = this.result.RowsAffected()
	return int(n)
}

func (this *Result) LastInsertId() (int64, error) {
	return this.result.LastInsertId()

}
func (this *Result) RowsAffected() (int64, error) {
	return this.result.RowsAffected()
}

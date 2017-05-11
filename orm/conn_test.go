package meorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestConn(t *testing.T) {
	Conn := NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai", "test")
	t.Log(Conn.dialect == 1)
	builder := Conn.NewBuilder().InsertInto("user").Columns("name").Values("name1")
	res := builder.Exec()
	t.Log(res.LastInsertId())
	t.Log(res.Err() == nil)
}

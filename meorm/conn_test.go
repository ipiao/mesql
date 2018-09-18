package meorm

import (
	"testing"

	"time"

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

func TestConnTx(t *testing.T) {
	Conn := NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai", "test")
	t.Log(Conn.dialect == 1)
	builder, err := Conn.BeginBuilder()
	t.Log(err == nil)
	sb1 := builder.InsertInto("user").Columns("name").Values("name1")
	res := sb1.Exec()
	t.Log(res.LastInsertId())
	t.Log(res.Err() == nil)

	sb2 := builder.InsertInto("user").Columns("name1").Values("name2")
	sb2.Exec()

	builder.Rollback()
}

func TestConnMultiTx(t *testing.T) {
	time1 := time.Now().UnixNano()
	conn := NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai", "test")
	time2 := time.Now().UnixNano()
	t.Log("conn time:", time2-time1)
	t.Log(conn.dialect == 1)
	tx1, err := conn.BeginBuilder()
	t.Log(err == nil)
	sb1 := tx1.InsertInto("user").Columns("name").Values("name1")
	sb1.Exec()

	//go func(conn *Conn) {
	tx2, err := conn.BeginBuilder()
	t.Log(err == nil)
	sb11 := tx2.InsertInto("user").Columns("name").Values("name3")
	sb11.Exec()
	sb12 := tx2.InsertInto("user").Columns("name").Values("name4")
	sb12.Exec()
	tx2.Commit()
	//}(conn)

	sb2 := tx1.InsertInto("user").Columns("name").Values("name2")
	sb2.Exec()

	tx1.Commit()
}

func BenchMarkBUilder(b *testing.B) {
	time1 := time.Now().UnixNano()
	conn := NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai", "test")
	time2 := time.Now().UnixNano()
	b.Log("conn time:", time2-time1)

	for i := 0; i < b.N; i++ {
		conn.NewBuilder().InsertInto("user").Columns("name").Values("nameN").Exec()
	}
}

func TestBuilderInter(t *testing.T) {
	time1 := time.Now().UnixNano()
	conn := NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai", "test")
	time2 := time.Now().UnixNano()
	t.Log("conn time:", time2-time1)

	res := conn.NewBuilder().Exec(nil)
	t.Log(res)
}

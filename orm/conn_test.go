package meorm

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Account   string `db:"account"`
	Password  string `db:"password"`
	Mobile    string `db:"phone"`
	Create_by int    `db:"create_by"`
	Update_by int    `db:"update_by"`
}

func TestConn(t *testing.T) {
	//	var conn1 = NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	//	t.Log(conn1.db.GetDB())

	//	var conn2 = NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	//	t.Log(conn2.db.GetDB())
	var count int

	var basedb, _ = sql.Open("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")

	var conn3 = MountConnection(basedb)
	conn3.db.Init()
	t.Log(basedb)

	t.Log(conn3.Name())
	go func() {
		//runtime.Gosched()
		var bb = conn3.Select("count(0)").From("user").Where("1=?", 1)
		bb.QueryNext(&count)
		t.Log(bb.ToSQL())
	}()

	var users []User
	var b = conn3.Select("name").From("user").WhereLikeL("name", "名字").WhereIn("id", 14, 15).Limit(10).Offset(0)
	t.Log(b.ToSQL())
	var count4, sql, err4 = b.CountCond()
	t.Log(count4, sql, err4)
	var n, err = b.QueryTo(&users)
	t.Log(users, n, err)
	//time.Sleep(1)
	t.Log(count)

	n, err = conn3.SQL("select name,account,id,phone from user where 1=1").QueryTo(&users)
	t.Log(users, n, err)
	//}()
}

func BenchmarkConn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var count int

		var basedb, _ = sql.Open("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")

		var conn3 = MountConnection(basedb)
		connections[conn3.name] = conn3
		b.Log(basedb)

		b.Log(conn3.Name())
		go func() {
			//runtime.Gosched()
			var bb = conn3.Select("count(0)").From("user").Where("1=?", 1)
			bb.QueryNext(&count)
			b.Log(bb.ToSQL())
		}()

		var users []User
		var bbb = conn3.Select("name").From("user").WhereLikeL("name", "名字").WhereIn("id", 14, 15).Limit(10).Offset(0)
		b.Log(bbb.ToSQL())
		var count4, sql, err4 = bbb.CountCond()
		b.Log(count4, sql, err4)
		var n, err = bbb.QueryTo(&users)
		b.Log(users, n, err)
		//time.Sleep(1)
		b.Log(count)

		//}()
	}

}

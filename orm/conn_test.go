package meorm

import (
	"database/sql"
	"testing"
	"time"

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

	r := conn3.Update("user").Set("name", "hahahabb").Where("1=?", 1).OrderBy("id desc").Limit(1)
	t.Log(r.ToSQL())
	r.Exec()

	ins := conn3.InsertInto("user").Columns("name", "account").Values("name1", "ac1", "das").Values("name2", "ac2")
	t.Log(ins.ToSQL())
	res := ins.Exec()
	t.Log(res.Err)

	n, err = conn3.SQL("select name,account,id,phone from user where 1=1").QueryTo(&users)
	t.Log(users, n, err)

	//}()

	//Delete 不测试，不建议物理删除
}

func BenchmarkConn(b *testing.B) {
	for i := 0; i < 30; i++ {
		var count int

		var basedb, _ = sql.Open("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")

		var conn3 = MountConnection(basedb)
		conn3.db.SetMaxOpenConns(300)
		b.Log(basedb)
		b.Log(time.Now())

		b.Log(conn3.Name())
		go func() {
			//runtime.Gosched()
			var bb = conn3.Select("count(0)").From("user").Where("1=?", 1)
			bb.QueryNext(&count)
			b.Log(bb.ToSQL())
		}()

		var users []User
		var sel = conn3.Select("name").From("user").WhereLikeL("name", "名字").WhereIn("id", 14, 15).Limit(10).Offset(0)
		b.Log(sel.ToSQL())
		var count4, sql, err4 = sel.CountCond()
		b.Log(count4, sql, err4)
		var n, err = sel.QueryTo(&users)
		b.Log(users, n, err)
		//time.Sleep(1)
		b.Log(count)

		r := conn3.Update("user").Set("name", "hahahabb").Where("1=?", 1).OrderBy("id desc").Limit(1)
		b.Log(r.ToSQL())
		r.Exec()

		ins := conn3.InsertInto("user").Columns("name", "account").Values("name1", "ac1").Values("name2", "ac2")
		b.Log(ins.ToSQL())
		ins.Exec()

		n, err = conn3.SQL("select name,account,id,phone from user where 1=1").QueryTo(&users)
		b.Log(time.Now())
		b.Log(users, n, err)

		//}()
	}

}

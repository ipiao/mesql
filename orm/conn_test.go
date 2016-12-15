package meorm

import (
	"database/sql"
	"fmt"
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
	Status    int
}

func TestConn(t *testing.T) {
	var datasource = "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai"
	var Conn = NewConnection("mysql", datasource)
	fmt.Println("conn", Conn, "db", Conn.DB())
	var users []User
	var sel = Conn.Select("name").From("user").WhereLikeL("name", "名字").WhereIn("id", []int{14, 15}).Limit(10).Offset(0)
	t.Log(sel.ToSQL())
	var count4, sql, err4 = sel.CountCond()
	t.Log(count4, sql, err4)
	var n, err = sel.QueryTo(&users)
	t.Log(users, n, err)
	//time.Sleep(1)
	t.Log(count4)

}

func BenchmarkConn(b *testing.B) {
	for i := 0; i < 30; i++ {
		var count int

		var basedb, _ = sql.Open("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")

		var conn3 = MountConnection(basedb)
		conn3.db.SetMaxOpenConns(3000)
		b.Log(basedb)
		b.Log(time.Now())

		b.Log(conn3.Name())
		//		go func() {
		//			//runtime.Gosched()
		//			conn3.Select("count(0)").From("user").Where("1=?", 1).QueryNext(&count)
		//			//b.Log(bb.ToSQL())
		//		}()

		var users []User
		var sel = conn3.Select("name").From("user").WhereLikeL("name", "名字").WhereIn("id", []int{14, 15}).Limit(10).Offset(0)
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

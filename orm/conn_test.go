package meorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Account  string `db:"account"`
	Password string `db:"password"`
	Mobile   string `db:"phone"`
	CreateBy int    `db:"create_by"`
	UpdateBy int    `db:"update_by"`
	Status   int
}

func TestConn(t *testing.T) {
	var datasource = "ipiao:1001@tcp(192.168.1.201:3306)/web_from_pg?charset=utf8mb4&loc=Asia%2fShanghai"
	var Conn = NewConnection("mysql", datasource, "web_from_pg")
	var users []User
	var sel = Conn.Select("*").From("consignor_user").WhereIn("id", 1, "149").Limit(10).Offset(0)
	t.Log(sel.ToSQL())
	var count4, sql, err4 = sel.CountCond()
	t.Log(count4, sql, err4)
	var n, err = sel.QueryTo(&users)
	t.Log(users, n, err)
	//time.Sleep(1)
	t.Log(count4)

}

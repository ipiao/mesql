package meorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Account  string `db:"_"`
	Password string `db:"pwd"`
	Mobile   string `db:"phone"`
	CreateBy int    `db:"_"`
	UpdateBy int    `db:"_"`
	Status   int
}

func (User) TableName() string {
	return "consignor_user"
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
	// var u = User{
	// 	Name:   "从结构体插入011",
	// 	Mobile: "13777462204",
	// }
	// var us = []*User{&u}
	// builder := Conn.InsertInto("consignor_user").Columns("phone", "name").Models(us)
	// err5 := builder.Exec().Err()
	// t.Log(err5 == nil, err5)
	// t.Log(builder.ToSQL())
	builder2 := Conn.Update("consignor_user").Set("name", "hahahha").Where("id=?", 176)
	err6 := builder2.Exec().Err()
	t.Log(err6 == nil, err6)
	t.Log(builder2.ToSQL())
}

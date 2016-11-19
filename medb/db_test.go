package medb

import (
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

func TestDB(t *testing.T) {
	var users = []User{}
	var dbName = "test"
	var err1 = RegisterDB(dbName, "mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	t.Log("[err1]:", err1)
	var db = OpenDB(dbName)
	var rows = db.Query(`select * from user`)
	//var err2 = rows.Close()
	//t.Log("[err2]:", err2)
	var cols = rows.Columns()
	t.Log("[columns]:", cols)
	var _, err3 = rows.ScanTo(&users)
	t.Log("[err3]:", err3, "[users]:", users)

	var stmt = db.Prepare(`select id,name,account,password,phone from user`)
	t.Log("[stmt]:", stmt)
	var n, err4 = stmt.Query().ScanTo(&users)
	t.Log("[err4]:", err4, "[users]:", users, "[n]:", n)
	stmt.Close()
	var n5, err5 = stmt.Query().ScanTo(&users)
	t.Log("[err5]:", err5, "[users]:", users, "[n5]:", n5)
	var n6, err6 = db.Exec(`update user set name = "asfsaagfhydadf"`).RowsAffected()
	t.Log("[n]:", n6, err6)
}

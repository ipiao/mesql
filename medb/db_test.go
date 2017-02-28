package medb

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Tt
	Consignor Consignor
}

type Tt struct {
	CreateTime string
	UpdateTime time.Time
}

type Consignor struct {
	ConsignorCode string
}

func TestDB(t *testing.T) {
	var users = []User{}
	var dbName = "test"
	var err1 = RegisterDB(dbName, "mysql", "ipiao:1001@tcp(192.168.1.201:3306)/hxz-web?charset=utf8mb4&loc=Asia%2fShanghai")
	t.Log("[err1]:", err1)
	var db = OpenDB(dbName)
	var rows = db.Query(`select * from consignor_user`)
	var _, err3 = rows.ScanTo(&users)
	t.Log("[err3]:", err3)
	t.Log(users)
}

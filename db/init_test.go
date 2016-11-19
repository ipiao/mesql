package medb

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestInit(t *testing.T) {
	var dbName = "test"
	var err1 = RegisterDB(dbName, "mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	t.Log("[medb]:", err1)
	var err2 = RegisterDB(dbName, "mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	t.Log("[medb]:", err2)
	var db = OpenDB(dbName)
	t.Log(db)
}

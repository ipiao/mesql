package medb

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMeDB(t *testing.T) {
	var datasource = "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai"
	err := RegisterDB("test", "mysql", datasource)
	if err != nil {
		t.Fatal(err)
	}
	db := OpenDB("test")
	err = db.Ping()
    t.
}

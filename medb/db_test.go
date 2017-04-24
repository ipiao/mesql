package medb

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMeDB(t *testing.T) {
	var datasource = "ipiao:1001@tcp(192.168.1.201:3306)/web_from_pg?charset=utf8mb4&loc=Asia%2fShanghai"
	err := RegisterDB("web_from_pg", "mysql", datasource)
	if err != nil {
		t.Fatal(err)
	}
	db := OpenDB("web_from_pg")
	err = db.Ping()
	t.Log(err == nil)
	for i := 0; i < 100; i++ {
		rows := db.Prepare("select * from consignor_user").Query()
		defer rows.Close()
	}
}

func BenchmarkMeDB(b *testing.B) {
	var datasource = "ipiao:1001@tcp(192.168.1.201:3306)/web_from_pg?charset=utf8mb4&loc=Asia%2fShanghai"
	err := RegisterDB("web_from_pg", "mysql", datasource)
	if err != nil {
		b.Fatal(err)
	}
	db := OpenDB("web_from_pg")
	err = db.Ping()
	b.Log(err == nil)

	for i := 0; i < b.N; i++ {
		rows := db.Prepare("select * from consignor_user").Query()
		defer rows.Close()
	}
}

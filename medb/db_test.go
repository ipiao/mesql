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
	for i := 0; i < 10; i++ {
		//rows := db.Prepare("select * from consignor_user").Query()
		rows := db.Query("select * from consignor_user")
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

func TestDBContxet(t *testing.T) {
	var datasource = "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai"
	err := RegisterDB("test", "mysql", datasource)
	if err != nil {
		t.Fatal(err)
	}
	db := OpenDB("test")
	err = db.Ping()
	t.Log(err == nil)

	row := db.Query("select * from user limit 1")
	cols, err := row.Columns()
	t.Log("cols:", cols, err == nil)
	cts, err := row.ColumnTypes()
	t.Log("cts:", cts, err == nil)
	for _, ct := range cts {
		dtName := ct.DatabaseTypeName()
		name := ct.Name()
		length, ok := ct.Length()
		precision, scale, bo := ct.DecimalSize()
		t.Log(dtName, ":", name, length, ok, precision, scale, bo)
	}
}

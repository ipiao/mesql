package medb

import (
	"context"
	"database/sql"
	"testing"

	"time"

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
	stats := db.Stats()
	t.Log(stats)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = db.PingContext(ctx)
	t.Log(err == nil)
	// time.Sleep(time.Second * 3)
	// err = db.PingContext(ctx)
	// t.Log(err != nil, err)
	var opts = new(sql.TxOptions)
	err = db.BeginTx(ctx, opts)
	t.Log(err == nil)
	stmt := db.Prepare(`insert into user(name,data,age,date)values('name1','data1',1,now())`)
	defer stmt.Close()
	t.Log(err == nil, err)
	res := stmt.Exec()
	t.Log(res.Err() == nil, res.Err())
	t.Log(res.RowsAffected())
	// time.Sleep(time.Second * 3)
	err = db.Commit()
	t.Log(err == nil, err)
}

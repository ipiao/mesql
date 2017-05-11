package medb

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestDB(t *testing.T) {
	err := RegisterDB("test", "mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai")
	if err != nil {
		t.Fatal("register db err:", err)
	}
	db := OpenDB("test")
	tx1, err := db.Begin()
	if err != nil {
		t.Fatal("tx1 begin err:", err)
	}
	tx1.Exec("insert into user(name)values('tx1')").Err()

	go func(db *DB) {
		tx2, err := db.Begin()
		if err != nil {
			t.Fatal("tx2 begin err:", err)
		}
		tx2.Exec("insert into user(name)values('tx2')").Err()

		tx2.Commit()
	}(db)

	err = tx1.Exec("insert into user(name1)values('tx3')").Err()

	err = tx1.Exec("insert into user(name)values('tx4')").Err()

	time.Sleep(time.Second)
	err = tx1.Commit()

	t.Fatal("commit  err", err)
}

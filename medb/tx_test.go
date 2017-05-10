package medb

import "testing"
import "time"
import "database/sql"

func TestTX(t *testing.T) {
	//
	var datasource = "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai"
	err := RegisterDB("test", "mysql", datasource)
	if err != nil {
		t.Fatal(err)
	}
	db := OpenDB("test")
	//
	err = db.Begin()
	t.Log("tx1 begin")

	err = db.Exec("insert into user(name)values('tx4')").Err()
	t.Log(19, db.autoCommit)
	if err != nil {
		db.Rollback()
		t.Fatal("执行失败", err)
	}

	// 模拟另一个执行方法
	go func(db2 *DB) {
		t.Log("tx2 begin")
		db2.Begin()
		t.Log(29, db.autoCommit)
		if err != nil {
			t.Fatal("事务开启失败", err)
		}
		err = db2.Exec("insert into user(name)values('tx2')").Err()
		if err != nil {
			t.Fatal("执行失败", err)
		}
		err = db2.Commit()
		t.Log(38, db.autoCommit)
		t.Log("tx2 commit")
		if err != nil {
			t.Fatal("提交失败", err)
		}
	}(db)

	err = db.Exec("insert into user(name)values('tx3')").Err()
	if err != nil {
		db.Rollback()
		t.Fatal("执行失败", err)
	}

	time.Sleep(time.Second * 1)
	if err != nil {
		t.Fatal("事务开启失败", err)
	}

	err = db.Exec("insert into user(name)values('tx5')").Err()
	if err != nil {
		db.Rollback()
		t.Fatal("执行失败", err)
	}

	err = db.Exec("insert into user(name)values('tx1')").Err()
	if err != nil {
		db.Rollback()
		t.Fatal("执行失败", err)
	}
	t.Log(60, db.autoCommit)
	err = db.Commit()
	t.Log("tx1 commit")
	if err != nil {
		t.Fatal("提交失败", err)
	}
}

func TestTXProp(t *testing.T) {
	//
	var datasource = "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai"
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		t.Fatal(err)
	}
	//
	tx1, err := db.Begin()
	t.Log(tx1, "tx1 begin")

	_, err = tx1.Exec("insert into user(name)values('tx4')")
	if err != nil {
		tx1.Rollback()
		t.Fatal("执行失败", err)
	}

	// 模拟另一个执行方法
	go func(db2 *sql.DB) {
		tx2, err := db2.Begin()
		t.Log("tx2 begin")

		if err != nil {
			t.Fatal("事务开启失败", err)
		}
		_, err = tx2.Exec("insert into user(name)values('tx2')")
		if err != nil {
			t.Fatal("执行失败", err)
		}
		err = tx2.Commit()
		t.Log("tx2 commit")
		if err != nil {
			t.Fatal("提交失败", err)
		}
	}(db)

	_, err = tx1.Exec("insert into user(name)values('tx3')")
	if err != nil {
		tx1.Rollback()
		t.Fatal("执行失败", err)
	}

	time.Sleep(time.Second * 1)
	if err != nil {
		t.Fatal("事务开启失败", err)
	}

	_, err = tx1.Exec("insert into user(name)values('tx5')")
	if err != nil {
		tx1.Rollback()
		t.Fatal("执行失败", err)
	}

	_, err = tx1.Exec("insert into user(name1)values('tx1')")
	if err != nil {
		tx1.Rollback()
		t.Fatal("执行失败", err)
	}
	err = tx1.Commit()
	t.Log("tx1 commit")
	if err != nil {
		t.Fatal("提交失败", err)
	}
}

package medb_test

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/metools/datetool"
)

func TestDB(t *testing.T) {
	err := medb.RegisterDB("test", "mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai")
	if err != nil {
		t.Fatal("register db err:", err)
	}
	db := medb.OpenDB("test")
	tx1, err := db.Begin()
	if err != nil {
		t.Fatal("tx1 begin err:", err)
	}
	tx1.Exec("insert into user(name)values('tx1')").Err()

	go func(db *medb.DB) {
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

func TestScan(t *testing.T) {
	type Date struct {
		D string
	}
	type User struct {
		Name string
		Data *Date
		Age  int
		Date time.Time
	}
	err := medb.RegisterDB("test", "mysql", "root:1001@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai")
	if err != nil {
		t.Fatal("register db err:", err)
	}
	db := medb.OpenDB("test")
	var us []User
	_, err = db.Query(`select * from user`).ScanTo(us)
	if err != nil {
		t.Log(err)
	}
	for _, u := range us {
		t.Log(u)
		t.Log(&us)
		t.Log(&u)
	}
}

func TestScanXXX(t *testing.T) {
	err := medb.RegisterDB("test", "mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai")
	if err != nil {
		t.Fatal("register db err:", err)
	}
	db := medb.OpenDB("test")
	defer db.Close()

	// medb.Logger.Skip(3)
	medb.Logger.Skip(5)
	medb.Logger.SetLevel(4)

	medb.RegisterTimeParserFunc(datetool.ParseTime)

	// type Ager int

	type User struct {
		Id       int
		Name     string
		Age      int        `db:"col:age;"`
		CreateAt *time.Time `db:"col:created_at;"`
	}

	u := []User{}

	ret, err := db.Query("select * from user ").ScanTo(&u)
	if err != nil {
		t.Fatal("query db err:", err)
	}
	t.Log(ret)
	t.Log(u)

}

type Time time.Time

func (t *Time) MedbParse(s string) {
	tim, _ := datetool.ParseTime(s)
	*t = Time(tim)
}

func (t *Time) String() string {
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

func TestScanCp(t *testing.T) {
	err := medb.RegisterDB("test", "mysql", "root:justfortest@tcp(118.25.7.38:3306)/test?charset=utf8mb4&loc=Asia%2fShanghai")
	if err != nil {
		t.Fatal("register db err:", err)
	}
	db := medb.OpenDB("test")
	defer db.Close()

	medb.Logger.Skip(5)
	// medb.Logger.Skip(5)
	medb.Logger.SetLevel(1)

	medb.RegisterTimeParserFunc(datetool.ParseTime)

	type User struct {
		Id       int    `db:"id"`
		Name     string `db:"name"`
		Age      int    `db:"age"`
		CreateAt Time   `db:"col:create_at,cp"`
	}

	u := []User{}

	ret, err := db.Query("select * from user ").ScanTo(&u)
	if err != nil {
		t.Fatal("query db err:", err)
	}
	t.Log(ret)
	t.Log(u)

}

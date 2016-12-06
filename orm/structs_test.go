package meorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type Info struct {
	Role int
}

func (this User) TableName() string {
	this.Name = "user"
	return this.Name
}

func TestStructIt(t *testing.T) {
	var conn = NewConnection("mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
	var u = User{Id: 2, Name: "hello", Status: 30}
	var u2 = User{Id: 1, Name: "hello2"}
	var users = []User{u, u2}
	var res = conn.InsertModels(&users)
	t.Log(res.Err)
}

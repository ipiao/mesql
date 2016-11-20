package mesql

//"depot/models"

type User struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Account   string `db:"account"`
	Password  string `db:"password"`
	Mobile    string `db:"phone"`
	Create_by int    `db:"create_by"`
	Update_by int    `db:"update_by"`
}

//
//func TestMesql(t *testing.T) {
//	//RegisterDriver("mysql",driver.)
//	var users []User
//	var err = RegisterDB("test", "mysql", "root:1001@tcp(127.0.0.1:3306)/depot?charset=utf8mb4&loc=Asia%2fShanghai")
//	var db = OpenDB("test")
//	defer db.Close()
//	conn := runner.NewConnection(db.db, "mysql")
//	t.Log(conn)
//	_, err = conn.Update("user").Set("name", "hello").Where("id=?", 14).Exec()
//	t.Log(users)
//	err = conn.SQL("select id,name,update_by,create_by,phone from user where id=?", 14).QueryStructs(&users)
//	t.Log(users, err)
//	//var u = []User{{Name: "ykk2", Account: "ykkk", Password: "123456"}, {Name: "ykk4"}}
//	//n, err := db.InsertModels(&u)
//	//t.Log(n, err)
//	t.Log(err, dbs, db)
//	var r = db.Ping()
//	t.Log(r)
//	var rows = db.Query("select id,name,update_by,create_by,phone from `user` ")
//	var id int
//	var name string
//	t.Log(rows.Error())
//	cols, err := rows.Columns()
//	t.Log(cols, err)

//	//	if rows.Next() {
//	//	}
//	//err = rows.ScanNext(&id, &name)
//	rows.ScanTo(&users)
//	t.Log(users, id, name, err)
//}

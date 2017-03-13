package medb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

// 自定义DB，包含db连接信息
type DB struct {
	name       string
	db         *sql.DB
	tx         *sql.Tx
	autocommit bool
}

// 初始化db
func (this *DB) Init() {
	this.db.SetMaxOpenConns(DefaultMaxOpenConns)
	this.db.SetMaxIdleConns(DefaultMaxIdleConns)
	this.db.SetConnMaxLifetime(DefaultConnMaxLifetime)
}

func (this *DB) Name() string {
	return this.name
}

// 嵌入db
func (this *DB) MountDB(db *sql.DB) error {
	mu.Lock()
	defer mu.Unlock()
	if this.db != nil {
		return ErrMountDB
	}
	var defaultname = RandomName()
	this.db = db
	this.name = defaultname
	this.autocommit = true
	dbs[defaultname] = this
	return nil
}

// 设置最大连接数
func (this *DB) SetMaxOpenConns(n int) {
	this.db.SetMaxOpenConns(n)
}

// 设置最大空闲连接数
func (this *DB) SetMaxIdleConns(n int) {
	this.db.SetMaxIdleConns(n)
}

// 设置连接时间
func (this *DB) SetConnMaxLifetime(d time.Duration) {
	this.db.SetConnMaxLifetime(d)
}

// DB的驱动
func (this *DB) Driver() driver.Driver {
	return this.db.Driver()
}

// Ping检查连接是否有效
func (this *DB) Ping() error {
	return this.db.Ping()
}

// Close,关闭连接
func (this *DB) Close() error {
	return this.db.Close()
}

// 解析sql
func (this *DB) Exec(sql string, params ...interface{}) *Result {
	// log.Println("[MESQL]:", sql, "[ARGS]:", params)
	if this.autocommit {
		var res, err = this.db.Exec(sql, params...)
		return &Result{result: res, Err: err}
	} else {
		var res, err = this.tx.Exec(sql, params...)
		return &Result{result: res, Err: err}
	}
}

// Call 调用存储过程
func (this *DB) Call(procedure string, params ...interface{}) *Rows {
	var sql = "call " + procedure + "("
	for i := 0; i < len(params); i++ {
		if i > 0 {
			sql += " ,?"
		} else {
			sql += "?"
		}
	}
	sql += ")"
	var rows, err = this.db.Query(sql, params...)
	return &Rows{rows, err, nil}
}

// 查询
func (this *DB) Query(sql string, params ...interface{}) *Rows {
	// log.Println("[MESQL]:", sql, "[ARGS]:", params)
	if this.autocommit {
		var rows, err = this.db.Query(sql, params...)
		return &Rows{rows: rows, err: err}

	} else {
		var rows, err = this.tx.Query(sql, params...)
		return &Rows{rows: rows, err: err}
	}

}

// 查询最多单行
func (this *DB) QueryRow(sql string, params ...interface{}) *Row {
	if this.autocommit {
		var row = this.db.QueryRow(sql, params...)
		return &Row{row: row}
	} else {
		var row = this.tx.QueryRow(sql, params...)
		return &Row{row: row}
	}
}

// 预处理
func (this *DB) Prepare(sql string) *Stmt {
	if this.autocommit {
		var stmt, err = this.db.Prepare(sql)
		return &Stmt{stmt: stmt, err: err}
	} else {
		var stmt, err = this.tx.Prepare(sql)
		return &Stmt{stmt: stmt, err: err}
	}
}

// 开启事务
func (this *DB) Begin() bool {
	var err error
	if this.autocommit {
		this.tx, err = this.db.Begin()
		if err != nil {
			return false
		}
		this.autocommit = false
	}
	return true
}

// 提交事务
func (this *DB) Commit() error {
	if this.autocommit {
		return ErrCommit
	}
	var err = this.tx.Commit()
	this.tx = nil
	this.autocommit = true
	return err
}

// 回滚
func (this *DB) RollBack() error {
	if this.autocommit {
		return ErrRollBack
	}
	var err = this.tx.Rollback()
	this.tx = nil
	this.autocommit = true
	return err
}

// 使已经存在的状态生成事务的状态
func (this *DB) Stmt(stmt *Stmt) *Stmt {
	if this.autocommit {
		return &Stmt{err: errors.New("事务没有开启")}
	}
	var s = this.tx.Stmt(stmt.stmt)
	return &Stmt{stmt: s}
}

// 返回连接的信息
func (this *DB) Stats() sql.DBStats {
	return this.db.Stats()
}

// 连接数
func (this *DB) OpenConnetcions() int {
	return this.db.Stats().OpenConnections
}

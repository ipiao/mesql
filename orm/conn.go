package meorm

import (
	"database/sql"
	"sync"

	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/orm/dialect"
	"github.com/mgutz/dat/common"
)

var (
	bufPool = common.NewBufferPool()
	mutex   = new(sync.Mutex)
)

const (
	ormTag            = "db"
	ormFieldSelectTag = medb.MedbFieldName
)

// Conn connection
// maybe it will have a pool
type Conn struct {
	*medb.DB
	dialect dialect.Dialect
}

// ShowSQL 打印日志
func ShowSQL(b bool) {
	medb.ShowSQL(b)
}

// MountConnection 直接移植已有数据连接
func MountConnection(name string, basedb *sql.DB, dialect dialect.Dialect) *Conn {
	mutex.Lock()
	defer mutex.Unlock()
	var medb = new(medb.DB)
	var err = medb.MountDB(basedb, name)
	if err != nil {
		panic(err)
	}
	var conn = &Conn{
		DB:      medb,
		dialect: dialect,
	}
	return conn
}

// NewConnection 新建连接
func NewConnection(driverName, dataSource string, name string) *Conn {
	mutex.Lock()
	defer mutex.Unlock()
	var err = medb.RegisterDB(name, driverName, dataSource)
	if err != nil {
		panic(err)
	}
	var conn = &Conn{
		DB:      medb.OpenDB(name),
		dialect: dialect.ConvertDriverNameToDialect(driverName),
	}
	return conn
}

// NewBuilder 构建builder
func (c *Conn) NewBuilder() *Builder {
	return &Builder{
		Executor: c.DB,
		dialect:  c.dialect,
	}
}

// BeginBuilder 事务
func (c *Conn) BeginBuilder() (*Builder, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Builder{
		Executor: tx,
		dialect:  c.dialect,
	}, nil
}

// SQL 直接写sql
func (c *Conn) SQL(sql string, args ...interface{}) *CommonBuilder {
	return c.NewBuilder().SQL(sql, args...)
}

// Select 生成查询构造器
func (c *Conn) Select(cols ...string) *SelectBuilder {
	return c.NewBuilder().Select(cols...)
}

// Update 生成更新构造器
func (c *Conn) Update(table string) *UpdateBuilder {
	return c.NewBuilder().Update(table)
}

// InsertOrUpdate 生成插入或更新构造器
func (c *Conn) InsertOrUpdate(table string) *InsupBuilder {
	return c.NewBuilder().InsertOrUpdate(table)
}

// InsertInto 生成插入构造器
func (c *Conn) InsertInto(table string) *InsertBuilder {
	return c.NewBuilder().InsertInto(table)
}

// ReplaceInto 生成插入构造器
func (c *Conn) ReplaceInto(table string) *InsertBuilder {
	return c.NewBuilder().ReplaceInto(table)
}

// DeleteFrom 生成删除构造器
func (c *Conn) DeleteFrom(table string) *DeleteBuilder {
	return c.NewBuilder().DeleteFrom(table)
}

// Delete 生成删除构造器
func (c *Conn) Delete(column string) *DeleteBuilder {
	return c.NewBuilder().Delete(column)
}

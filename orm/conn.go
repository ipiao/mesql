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

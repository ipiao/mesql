package meorm

import (
	"database/sql"
	"sync"

	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/orm/common"
	"github.com/ipiao/mesql/orm/dialect"
)

var (
	bufPool     = common.NewBufferPool()
	connections = map[string]*Conn{}
	mutex       = new(sync.Mutex)
)

const (
	ormTag            = "db"
	ormFieldSelectTag = medb.MedbFieldName
)

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
	connections[name] = conn
	return connections[name]
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
	connections[name] = conn
	return conn
}

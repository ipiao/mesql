package meorm

import (
	"database/sql"
	"sync"

	"github.com/ipiao/mesql/medb"
	"github.com/ipiao/mesql/orm/common"
)

var (
	bufPool     = common.NewBufferPool()
	connections = map[string]*Conn{}
	mutex       = new(sync.Mutex)
)

const (
	ormTag      = "db"
	ormFieldTag = "col"
)

// MountConnection 直接移植已有数据连接
func MountConnection(basedb *sql.DB, name string) *Conn {
	mutex.Lock()
	defer mutex.Unlock()
	var medb = new(medb.DB)
	var err = medb.MountDB(basedb, name)
	if err != nil {
		panic(err)
	}
	var conn = &Conn{
		DB:   medb,
		name: name,
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
		DB:   medb.OpenDB(name),
		name: name,
	}
	connections[name] = conn
	return conn
}

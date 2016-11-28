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

// 直接移植已有数据连接
func MountConnection(basedb *sql.DB) *Conn {
	mutex.Lock()
	defer mutex.Unlock()
	var medb = new(medb.DB)
	var err = medb.MountDB(basedb)
	var name = medb.Name()
	if err != nil {
		meLog.Debug(err)
	}
	var conn = &Conn{
		db:   medb,
		name: name,
	}
	connections[name] = conn
	return connections[name]
}

// 新建连接
func NewConnection(driverName, dataSource string, connname ...string) *Conn {
	mutex.Lock()
	defer mutex.Unlock()
	var name string
	if len(connname) == 0 {
		name = medb.RandomName()
	} else {
		name = connname[0]
	}
	var err = medb.RegisterDB(name, driverName, dataSource)
	if err != nil {
		panic(err)
	}
	var conn = &Conn{
		db:   medb.OpenDB(name),
		name: name,
	}
	connections[name] = conn
	return conn
}

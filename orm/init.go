package meorm

import (
	"database/sql"
	"ipiao/mesql/medb"
	"ipiao/mesql/orm/common"
)

var bufPool = common.NewBufferPool()
var conn *Conn

func NewConnection(basedb *sql.DB) *Conn {
	var medb = medb.DB{}
	var err = medb.MountDB(basedb)
	if err != nil {
		meLog.Debug(err)
	}
	medb.Init()
	return &Conn{
		db: &medb,
	}
}

package pool

import "github.com/ipiao/mesql/medb"

var (
	defaultPrefix = "default_"
	masterPrefix  = "master_"
	slavePrefix   = "slave_"
)

// 连接池，针对的是同一个业务模块
// 比如说，主从分库
type DBManager struct {
	dbm map[string][]struct {
		Config *DBPoolConfig
		dbs    []*medb.DB
	}
}

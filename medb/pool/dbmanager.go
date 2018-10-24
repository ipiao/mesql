package pool

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ipiao/mesql/medb"
)

const (
	DBTypeDefault = "default"
	DBTypeMaster  = "master"
	DBTypeSlave   = "slave"
)

// 主要是分库，主从
type DBManagerConfig struct {
	Masters []*DBPoolConfig
	Slaves  []*DBPoolConfig
}

// 连接池，针对的是同一个业务模块
// 比如说，主从分库
type DBManager struct {
	selector Selector
	masters  map[string]*DBPool
	slaves   map[string]*DBPool
	conf     *DBManagerConfig
	putDBs   map[string]*DBPool
}

// 创建一个连接管理器
func NewDBManger(cfs []*DBPoolConfig) (m *DBManager, err error) {
	m = new(DBManager)
	m.putDBs = make(map[string]*DBPool)
	m.SetSelector(defaultSelector)
	conf := &DBManagerConfig{
		Masters: make([]*DBPoolConfig, 0),
		Slaves:  make([]*DBPoolConfig, 0),
	}
	for _, c := range cfs {
		if c.Type == DBTypeMaster {
			if c.Size == 0 {
				c.Size = 1 // 主库的连接池大小默认设置1
			}
			conf.Masters = append(conf.Masters, c)
		} else if c.Type == DBTypeSlave {
			conf.Slaves = append(conf.Slaves, c)
		} else {
			err = fmt.Errorf("wrong type of db:%s", c.Type)
			return
		}
	}
	m.conf = conf
	err = m.Connect()
	return
}

// 自连接
func (m *DBManager) Connect() error {
	if m.conf == nil {
		return errors.New("config can not be empty")
	}
	if len(m.masters) < len(m.conf.Masters) {
		m.masters = make(map[string]*DBPool)
		for _, c := range m.conf.Masters {
			pool, err := NewDBPool(c)
			if err != nil {
				return err
			}
			m.masters[c.Database] = pool
		}
	}

	if len(m.slaves) < len(m.conf.Slaves) {
		m.slaves = make(map[string]*DBPool)
		for _, c := range m.conf.Slaves {
			pool, err := NewDBPool(c)
			if err != nil {
				return err
			}
			m.slaves[c.Database] = pool
		}
	}
	return nil
}

type Selector func([]*DBPoolConfig) int

var defaultSelector = func(cfs []*DBPoolConfig) int {
	r := rand.Int()
	return r % len(cfs)
}

func (m *DBManager) SetSelector(sel Selector) {
	m.selector = sel
}

func (m *DBManager) GetMasterPool(sels ...Selector) *DBPool {
	sel := m.selector
	if len(sels) > 0 {
		sel = sels[0]
	}
	i := sel(m.conf.Masters)
	dbName := m.conf.Masters[i].Database
	return m.masters[dbName]
}

func (m *DBManager) GetMasterDB(sels ...Selector) (*medb.DB, error) {
	pool := m.GetMasterPool(sels...)
	db, err := pool.GetDB()
	if err != nil {
		return nil, err
	}
	if m.putDBs[db.Name()] == nil {
		m.putDBs[db.Name()] = pool
	}
	return db, nil
}

func (m *DBManager) GetSlavePool(sels ...Selector) *DBPool {
	sel := m.selector
	if len(sels) > 0 {
		sel = sels[0]
	}
	i := sel(m.conf.Slaves)
	dbName := m.conf.Slaves[i].Database
	return m.slaves[dbName]
}

func (m *DBManager) GetSlaveDB(sels ...Selector) (*medb.DB, error) {
	pool := m.GetSlavePool(sels...)
	db, err := pool.GetDB()
	if err != nil {
		return nil, err
	}
	if m.putDBs[db.Name()] == nil {
		m.putDBs[db.Name()] = pool
	}
	return db, nil
}

func (m *DBManager) PutDB(db *medb.DB) error {
	if pool, ok := m.putDBs[db.Name()]; ok {
		pool.PutDB(db)
		return nil
	}
	return fmt.Errorf("unkonwn db %s", db.Name())
}

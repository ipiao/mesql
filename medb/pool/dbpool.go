package pool

import (
	"context"
	"fmt"
	"time"

	"github.com/ipiao/mesql/medb"
)

// 创建dbpool的配置
// 定义好字段，免去调用方自己定义
type DBPoolConfig struct {
	User            string
	Password        string
	Host            string
	Port            string
	Database        string
	Loc             string // Asia%2fShanghai
	Charset         string // utf8mb4
	Uri             string
	Type            string //
	Size            int
	DriverName      string
	Monitor         bool // 监视
	MonitorInterval time.Duration
	GetDBTimeOut    time.Duration
}

const (
	DBTypeDefault = "default"
	DBTypeMaster  = "master"
	DBTypeSlaver  = "slave"
)

// SetDefault 设置默认值
func (c *DBPoolConfig) SetDefault() {
	if len(c.User) == 0 {
		c.User = "root"
	}
	if len(c.Host) == 0 {
		c.Host = "127.0.0.1"
	}
	if len(c.Port) == 0 {
		c.Port = "3306"
	}
	if len(c.Loc) == 0 {
		c.Loc = "Asia%2fShanghai"
	}
	if len(c.Charset) == 0 {
		c.Charset = "utf8mb4"
	}
	if len(c.Database) == 0 {
		c.Database = "test"
	}
	if len(c.Type) == 0 {
		c.Type = DBTypeDefault
	}
	if len(c.DriverName) == 0 {
		c.DriverName = "mysql"
	}
	if c.Size == 0 {
		c.Size = 5
	}
	if c.GetDBTimeOut == 0 {
		c.GetDBTimeOut = time.Millisecond * 100
	}
	if c.MonitorInterval == 0 {
		c.GetDBTimeOut = time.Minute * 10
	}
	_ = c.String()
}

func (c *DBPoolConfig) String() string {
	if len(c.Uri) == 0 {
		c.Uri = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s", c.User, c.Password, c.Host, c.Port, c.Database, c.Charset, c.Loc)
	}
	return fmt.Sprintf("%s:%s", c.Type, c.Uri)
}

type DBPool struct {
	Config      *DBPoolConfig
	dbs         chan *medb.DB
	lastGetTime time.Time
}

func (p *DBPool) GetDB() (db *medb.DB, err error) {
	return p.GetDBTimeOut(p.Config.GetDBTimeOut)
}

func (p *DBPool) GetDBTimeOut(timeOut time.Duration) (db *medb.DB, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()
	select {
	case db = <-p.dbs:
		p.lastGetTime = time.Now()
		return
	case <-ctx.Done():
		err = ctx.Err()
		return
	}
}

func (p *DBPool) PutDB(db *medb.DB) {
	if db != nil {
		p.dbs <- db
	}
}

func (p *DBPool) CreateNewDB(i int) (err error) {
	name := RandomDBName(i)
	err = medb.RegisterDB(name, p.Config.DriverName, p.Config.Uri)
	if err != nil {
		return
	}
	db := medb.OpenDB(name)
	p.PutDB(db)
	return
}

func NewDBPool(c *DBPoolConfig) (pool *DBPool, err error) {
	c.SetDefault()
	pool = &DBPool{Config: c, dbs: make(chan *medb.DB, c.Size)}
	for i := 0; i < c.Size; i++ {
		name := RandomDBName(i)
		err = medb.RegisterDB(name, "mysql", c.Uri)
		if err != nil {
			return
		}
		pool.PutDB(medb.OpenDB(name))
	}
	pool.RunMonitor()
	return
}

func RandomDBName(i int) string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), i)
}

func (p *DBPool) RunMonitor() {
	if p.Config.Monitor {
		go func() {
			ticker := time.NewTicker(p.Config.MonitorInterval)
			for {
				select {
				case <-ticker.C:
					medb.Logger.Infof("DBPool Monitor Task Start")

					if time.Now().Sub(p.lastGetTime) > time.Minute { // 上次获取时间到现在超过一分钟了，说明有连接没能归还
						num := p.Config.Size - len(p.dbs)
						for i := 0; i < num; i++ {
							p.CreateNewDB(i)
						}
					}
					if len(p.dbs) == p.Config.Size { // 是否处于空闲状态

						for i := 0; i < p.Config.Size; i++ {
							db := <-p.dbs
							if err := db.Ping(); err != nil {
								medb.Logger.Infof("Monitor Ping Error:", err)
								for i := 0; i < 3; i++ {
									err = p.CreateNewDB(i)
									if err == nil {
										break
									}
								}
							} else {
								p.PutDB(db)
							}
						}
					}
				}
			}
		}()
	}
}

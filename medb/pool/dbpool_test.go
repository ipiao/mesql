package pool

import (
	"os"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	_ "github.com/go-sql-driver/mysql"
)

func TestSetDefault(t *testing.T) {
	conf := new(DBPoolConfig)
	conf.Host = "118.25.7.38"
	conf.Password = "yukktop001"
	conf.SetDefault()
	conf.GetDBTimeOut = 0
	conf.MonitorInterval = time.Second * 2
	conf.Monitor = true
	t.Log(conf)

	_, err := NewDBPool(conf)
	assert.Equal(t, err, nil)

	// for i := 0; i < conf.Size; i++ {
	// 	db, err := pool.GetDB()
	// 	t.Log(db)
	// 	if err != nil {
	// 		t.Log(err)
	// 	}
	// 	go func(dbb *medb.DB) {
	// 		time.Sleep(time.Millisecond * 5)
	// 		pool.PutDB(db)
	// 	}(db)
	// }

	// db, err := pool.GetDB()
	// t.Log(db)
	// t.Log(len(pool.dbs))
	// if err != nil {
	// 	t.Log(err)
	// }
	// t.Log(conf)
	time.Sleep(time.Second * 10)
	os.Exit(0)
}

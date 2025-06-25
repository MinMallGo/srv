package initialize

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"srv/inventorySrv/global"
)

func InitRedLock() {
	pool := goredis.NewPool(goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", global.SrvConfig.Redis.Host, global.SrvConfig.Redis.Port),
		Password: global.SrvConfig.Redis.Password,
	}))

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	global.RedLock = redsync.New(pool)
}

package common

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var RDB *redis.Client

func InitRedis() error {
	conf := &config.AppConf.Redis
	if conf.Url == "" || conf.Port == 0 {
		return fmt.Errorf("redis conf error")
	}
	addr := conf.Url + ":" + strconv.FormatInt(int64(conf.Port), 10)
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis connect error: %s", err)
	}

	return nil
}

func init() {
	initializer.Register("redis", InitRedis)
}

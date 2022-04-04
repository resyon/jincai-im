package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/conf"
	"sync"
)

var (
	once  sync.Once
	_conn *redis.Conn
)

func GetRedis() *redis.Conn {
	once.Do(func() {
		client := NewRedisClient()
		_conn = client.Conn(context.TODO())

	})
	return _conn
}

func NewRedisClient() *redis.Client {
	options := conf.GetRedisConf()
	client := redis.NewClient(&options.Redis)
	return client
}

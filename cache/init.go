package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/conf"
	"sync"
)

var (
	once      sync.Once
	_conn     *redis.Conn
	RedisPool = &redisPool{userMap: &sync.Map{}}
)

type redisPool struct {
	userMap *sync.Map
}

func (p *redisPool) GetRedisConnection(userId int) *redis.Client {
	_c, ok := p.userMap.Load(userId)
	if !ok {
		c := NewRedisClient()
		p.userMap.Store(userId, _c)
		return c
	}
	return _c.(*redis.Client)
}

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

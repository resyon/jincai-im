package conf

import (
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/log"
	"github.com/spf13/viper"
	"sync"
)

const confType = "yaml"
const confPath = ".."

const redisConfName = "redis"

var (
	_redisConf RedisConf
	once       sync.Once
)

type RedisConf struct {
	//Addr     string `yaml:"addr"`
	//DB       int    `yaml:"db"`
	//Password string `yaml:"password"`
	Redis redis.Options
}

func GetRedisConf() RedisConf {
	once.Do(func() {
		viper.SetConfigType(confType)
		viper.SetConfigName(redisConfName)
		viper.AddConfigPath(confPath)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.LOG.Panicf("Fail to read redis_conf, err=%+v\n", err)
		}
		err := viper.Unmarshal(&_redisConf)
		if err != nil {
			log.LOG.Panicf("Fail to parse redis_conf, err=%+v\n", err)
		}
	})
	return _redisConf
}

package conf

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sync"
)

const redisConfType = "yaml"
const redisConfName = "redis"
const redisConfPath = ".."

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
		viper.SetConfigType(redisConfType)
		viper.SetConfigName(redisConfName)
		viper.AddConfigPath(redisConfPath)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		err := viper.Unmarshal(&_redisConf)
		if err != nil {
			panic(err)
		}
	})
	return _redisConf
}

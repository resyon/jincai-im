package conf

import (
	"github.com/resyon/jincai-im/log"
	"github.com/spf13/viper"
	"sync"
)

const (
	appConfName = "app"
)

var (
	_appConf AppConf
	appOnce  sync.Once
)

type AppConf struct {
	Port int     `yaml:"port"`
	JWT  JWTConf `yaml:"jwt"`
}

type JWTConf struct {
	AuthKey string `yaml:"auth_key"`
}

func GetAppConf() AppConf {
	appOnce.Do(func() {
		viper.SetConfigType(confType)
		viper.SetConfigName(appConfName)
		viper.AddConfigPath(confPath)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.LOG.Panicf("Fail to read app_conf, err=%+v\n", err)
		}
		err := viper.Unmarshal(&_appConf)
		if err != nil {
			log.LOG.Panicf("Fail to parse app_conf, err=%+v\n", err)
		}
	})
	return _appConf
}

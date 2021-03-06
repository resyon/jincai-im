package conf

import (
	"fmt"
	"github.com/resyon/jincai-im/log"
	"github.com/spf13/viper"
	"sync"
)

//
//username: root
//password: abc
//dbname: jincai
//host: resyon.io
//port: 3306

const (
	mysqlConfName = "mysql"
	mysqlDsnFmt   = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4"
)

var (
	_dsn       string
	_mysqlConf MysqlConf
	mysqlOnce  sync.Once
)

type MysqlConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func GetMysqlDSN() string {
	mysqlOnce.Do(func() {
		viper.SetConfigType(confType)
		viper.SetConfigName(mysqlConfName)
		viper.AddConfigPath(confPath)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.LOG.Panicf("fail to get mysql conf, err=%+v", err)
		}
		err := viper.Unmarshal(&_mysqlConf)
		if err != nil {
			//#user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
			log.LOG.Panicf("fail to parse mysql conf, err=%+v", err)
		}
		_dsn = fmt.Sprintf(mysqlDsnFmt, _mysqlConf.Username,
			_mysqlConf.Password, _mysqlConf.Host, _mysqlConf.Port, _mysqlConf.DBName)
	})

	return _dsn
}

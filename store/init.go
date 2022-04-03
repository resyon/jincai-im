package store

import (
	"github.com/resyon/jincai-im/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	_db  *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		dsn := conf.GetMysqlDSN()
		var err error
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
	return _db
}

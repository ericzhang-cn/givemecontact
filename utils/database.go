package utils

import (
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sync"
)

var db *gorm.DB
var onceDb sync.Once

func GetDB() *gorm.DB {
	onceDb.Do(func() {
		_db, err := gorm.Open(sqlite.Open(viper.GetString("sqlite.db_file")), &gorm.Config{})
		if err != nil {
			log.Fatalf("open database failed, error message: %s", err.Error())
		}
		db = _db
	})
	return db
}

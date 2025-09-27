package bootstrap

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"trip-map/config"
	"trip-map/global"
)

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // 日志级别: Silent, Error, Warn, Info
			Colorful:      true,         // 是否彩色打印
		},
	)

	dsn := config.SysConfig.Mysql.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(config.SysConfig.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.SysConfig.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.Db = db
}

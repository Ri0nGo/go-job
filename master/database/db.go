package database

import (
	"fmt"
	"go-job/master/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	once    sync.Once
	mysqlDb *gorm.DB
)

func NewMySQLWithGORM() *gorm.DB {
	once.Do(initMySQLDBWithGorm)
	return mysqlDb
}

func initMySQLDBWithGorm() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.App.MySQL.Username,
		config.App.MySQL.Password,
		config.App.MySQL.Host,
		config.App.MySQL.Port,
		config.App.MySQL.Database,
	)
	gromCfg := &gorm.Config{}
	if config.App.MySQL.ShowSQL {
		gormLog := logger.New(
			log.New(os.Stdout, "\r\n[GORM] ", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Info,            // Log level
				IgnoreRecordNotFoundError: true,                   // 忽略 ErrRecordNotFound
				Colorful:                  true,                   // 彩色打印
			},
		)
		gromCfg.Logger = gormLog
	}
	db, err := gorm.Open(mysql.Open(dsn), gromCfg)
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(config.App.MySQL.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.App.MySQL.MaxOpenConn)
	slog.Info("connect to mysql success")

	mysqlDb = db
}

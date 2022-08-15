package db

import (
	"fmt"
	"time"

	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/model"

	"github.com/zhfreal/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var (
		err  error
		dial gorm.Dialector
	)

	switch config.DB {
	case "mysql":
		dial = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Mysql.User,
			config.Mysql.Password,
			config.Mysql.Host,
			config.Mysql.Port,
			config.Mysql.DB,
		))
	case "sqlite":
		dial = sqlite.Open(config.Sqlite.DB)
	}

	if dial == nil {
		zap.S().Fatalw("failed to get dial, please check your config")
	}
	DB, err = gorm.Open(dial, &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		zap.S().Fatalw("failed to open db", "error", err)
	}
	DB.AutoMigrate(&model.Client{})
}

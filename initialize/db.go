package initialize

import (
	"EasyChatApp/global"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.ServiceConfig.DB.User,
		global.ServiceConfig.DB.Password,
		global.ServiceConfig.DB.Host,
		global.ServiceConfig.DB.Port,
		global.ServiceConfig.DB.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}

func InitRedis() {
	opt := redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.ServiceConfig.RedisDB.Host,
			global.ServiceConfig.RedisDB.Port),
		Password: "",
		DB:       10,
	}
	global.RedisDB = redis.NewClient(&opt)
}

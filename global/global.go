package global

import (
	"EasyChatApp/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ServiceConfig *config.ServiceConfig
	DB            *gorm.DB
	RedisDB       *redis.Client
)

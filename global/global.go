package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB
	MysqlDb    *gorm.DB
	PostgresDb *gorm.DB
	KingBaseDb *gorm.DB
	DmDb       *gorm.DB
	RedisDb    *redis.Client
)

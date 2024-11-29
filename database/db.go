package database

import (
	"context"
	"dbConnect-go-demo/config"
	"dbConnect-go-demo/database/drivers"
	"dbConnect-go-demo/database/models"
	"dbConnect-go-demo/global"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mu sync.Mutex

// InitDB 初始化数据库连接
func InitDB(cfg config.Config) error {
	// 初始化数据库连接 gorm
	for _, dbConfig := range cfg.Databases {
		fmt.Println(dbConfig)
		_, err := InitGormDB(dbConfig.Driver, dbConfig.DSN, dbConfig.MaxIdleConns, dbConfig.MaxOpenConns, dbConfig.ConnMaxLifetime)
		if err != nil {
			log.Fatalf("Failed to connect to database %s: %v", dbConfig.Driver, err)
		}
	}

	// 初始化 Redis 连接
	fmt.Println(cfg.Redis)
	if cfg.Redis.Addr != "" {
		_, err := InitRedisDB(cfg.Redis)
		if err != nil { // 初始化 Redis 连接失败
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
	}

	// 初始化其他连接（es、kafka、mongodb）等 todo

	return nil
}

// RedisClient 包含 Redis 客户端实例
type RedisClient struct {
	client *redis.Client
}

// InitRedisDB 初始化 Redis 客户端
func InitRedisDB(config config.RedisConfig) (*redis.Client, error) {

	if _, exists := GetManager().connections["redis"]; exists {
		return nil, nil
	}

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,     // Redis 地址
		Password: config.Password, // 密码
		DB:       config.DB,       // 数据库编号
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("Connected to Redis:", config.Addr)
	err = GetManager().AddConnection("redis", client)
	global.RedisDb = client // 设置全局 Redis 连接 todo 这两种方式都可以
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InitGormDB 封装具体的数据库连接逻辑
func InitGormDB(driver, dsn string, maxIdle, maxOpen, maxLifetime int) (interface{}, error) {
	mu.Lock()
	defer mu.Unlock()

	if db, exists := GetManager().connections[driver]; exists {
		return db, nil
	}

	var dialector gorm.Dialector
	switch driver {
	case "mysql": // mysql 数据库
		dialector = mysql.Open(dsn)
	case "postgres": // postgres 数据库
		dialector = postgres.Open(dsn)
	case "kingbase": // kingbase 数据库-人大金仓
		dialector = drivers.NewKingbaseDialector(dsn)
	case "dm": //达梦数据库
		dialector = drivers.NewDmDialector(dsn)
	default:
		return nil, errors.New("gorm Database connection not found")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//err = db.AutoMigrate(&models.User{}) // 自动迁移表结构 todo 改造成 go 协程处理
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdle)                                     //设置最大空闲连接数
	sqlDB.SetMaxOpenConns(maxOpen)                                     //设置可打开的最大连接数为 100 个
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second) //设置一个连接空闲后在多长时间内可复用

	log.Printf("gorm Database connected: %s", driver)

	// 添加到全局管理器
	err = GetManager().AddConnection(driver, db)
	SetGlobalGormDb(driver, db) // 设置全局数据库连接 todo 这两种方式都可以
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SetGlobalGormDb(driver string, db *gorm.DB) {
	switch driver {
	case "mysql": // mysql 数据库
		global.MysqlDb = db
	case "postgres": // postgres 数据库
		global.PostgresDb = db
	case "kingbase": // kingbase 数据库-人大金仓
		global.KingBaseDb = db
	case "dm": //达梦数据库
		global.DmDb = db
	default:
		log.Printf("gorm global config failed: %s", driver)
	}
}

func autoMigrate(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

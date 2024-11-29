package database

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"

	"gorm.io/gorm"
)

// DBManager 管理所有数据库连接
type DBManager struct {
	mu          sync.RWMutex
	connections map[string]interface{}
}

// 全局实例
var manager *DBManager
var once sync.Once

// GetManager 获取全局数据库管理器实例
func GetManager() *DBManager {
	once.Do(func() {
		manager = &DBManager{
			connections: make(map[string]interface{}),
		}
	})
	return manager
}

// AddConnection 添加连接实例
func (m *DBManager) AddConnection(name string, connection interface{}) error {
	if _, exists := m.connections[name]; exists {
		return fmt.Errorf("connection with name '%s' already exists", name)
	}
	m.connections[name] = connection
	return nil
}

// GetConnection 获取连接实例（通用）
func (m *DBManager) GetConnection(name string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("no connection found with name: %s", name)
	}
	return conn, nil
}

// GetGormDB 获取 GORM 数据库实例
func (m *DBManager) GetGormDB(name string) (*gorm.DB, error) {
	// 设置默认值
	if name == "" {
		name = "mysql"
	}
	conn, err := m.GetConnection(name)
	if err != nil {
		return nil, err
	}

	db, ok := conn.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("connection is not a GORM database instance: %s", name)
	}
	return db, nil
}

// GetRedisDB  获取 Redis 客户端实例
func (m *DBManager) GetRedisDB() (*redis.Client, error) {
	conn, err := m.GetConnection("redis")
	if err != nil {
		return nil, err
	}

	client, ok := conn.(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("connection is not a Redis client instance: redis")
	}
	return client, nil
}

// CloseConnections 关闭所有连接
func (m *DBManager) CloseConnections() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for name, conn := range m.connections {
		switch c := conn.(type) {
		case *gorm.DB: // 关闭 GORM 数据库连接
			sqlDB, err := c.DB()
			if err != nil {
				log.Printf("Failed to retrieve raw DB for connection '%s': %v", name, err)
				continue
			}
			err = sqlDB.Close()
			if err != nil {
				log.Printf("Failed to close GORM database connection '%s': %v", name, err)
			} else {
				log.Printf("Closed GORM database connection '%s'", name)
			}
		case *redis.Client: // 关闭 Redis 连接
			err := c.Close()
			if err != nil {
				log.Printf("Failed to close Redis connection '%s': %v", name, err)
			} else {
				log.Printf("Closed Redis connection '%s'", name)
			}
		default:
			log.Printf("Unknown connection type for '%s', skipping", name)
		}
	}
	m.connections = make(map[string]interface{})
	return nil
}

package drivers

import (
	"gorm.io/driver/postgres" // 金仓与 PostgreSQL 协议兼容
	"gorm.io/gorm"
)

// NewKingbaseDialector 金仓驱动
func NewKingbaseDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

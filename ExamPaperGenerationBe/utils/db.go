package utils

import (
	"graduation/mapper"

	"gorm.io/gorm"
)

// GetDB 返回全局数据库实例
func GetDB() *gorm.DB {
	return mapper.DB
}

package entity

import (
	"gorm.io/gorm"
	"time"
)

// User 表示用户实体
type User struct {
	ID                  int       `gorm:"primaryKey;column:id" json:"id"`
	Username            string    `gorm:"column:username" json:"username"`
	Password            string    `gorm:"column:password" json:"password"`
	UserRole            string    `gorm:"column:user_role" json:"user_role"`
	LastLogin           time.Time `gorm:"column:last_login" json:"last_login"`
	Enable              int       `gorm:"column:enable" json:"enable"`
	SimilarityThreshold float64   `gorm:"column:similarity_threshold" json:"similarity_threshold"`
}

// BeforeCreate 在创建用户记录前设置最后登录时间
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.LastLogin = time.Now()
	return nil
}

func (u *User) TableName() string {
	return "user" // 明确指定表名
}

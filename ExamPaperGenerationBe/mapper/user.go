package mapper

import (
	"fmt"
	"graduation/entity"

	"gorm.io/gorm"
)

// UserMapper 接口定义
type UserMapper struct {
	db *gorm.DB
}

// NewUserMapper 创建一个新的 UserMapper 实例
func NewUserMapper() *UserMapper {
	return &UserMapper{
		db: DB,
	}
}

// Login 用户登录
func (m *UserMapper) Login(user *entity.User) ([]entity.User, error) {
	var users []entity.User
	result := m.db.Where("username =? AND password =? AND enable = 1", user.Username, user.Password).Find(&users)
	return users, result.Error
}

// AddNewUser 添加新用户
func (m *UserMapper) AddNewUser(user *entity.User) (int64, error) {
	result := m.db.Create(user)
	return result.RowsAffected, result.Error
}

// UpdateLastLoginTime 更新用户最后登录时间
func (m *UserMapper) UpdateLastLoginTime(user *entity.User) error {
	result := m.db.Model(user).Where("username =? AND password =?", user.Username, user.Password).Update("last_login", user.LastLogin)
	return result.Error
}

// GetUserByUsername 根据用户名获取用户
func (m *UserMapper) GetUserByUsername(username string) ([]entity.User, error) {
	var users []entity.User
	result := m.db.Where("username =?", username).Find(&users)
	return users, result.Error
}

// GetApplyUser 获取待审核用户
func (m *UserMapper) GetApplyUser() ([]entity.User, error) {
	var users []entity.User
	result := m.db.Where("enable = 0 AND user_role = 'user'").Find(&users)
	return users, result.Error
}

// GetAllUser 获取所有用户
func (m *UserMapper) GetAllUser() ([]entity.User, error) {
	var users []entity.User
	result := m.db.Find(&users)
	return users, result.Error
}

// DeleteUser 删除用户
func (m *UserMapper) DeleteUser(username string) (int64, error) {
	result := m.db.Where("username =?", username).Delete(&entity.User{})
	return result.RowsAffected, result.Error
}

// PassApply 通过用户申请
func (m *UserMapper) PassApply(username string) (int64, error) {
	result := m.db.Model(&entity.User{}).Where("username =?", username).Update("enable", 1)
	return result.RowsAffected, result.Error
}

// DeleteApply 删除用户申请
func (m *UserMapper) DeleteApply(username string) (int64, error) {
	result := m.db.Where("username =?", username).Delete(&entity.User{})
	return result.RowsAffected, result.Error
}

// UpdateUserSimilarityThreshold 更新用户的相似度阈值
func UpdateUserSimilarityThreshold(username string, threshold float64) error {
	// 更新用户相似度阈值
	result := DB.Model(&entity.User{}).Where("username = ?", username).Update("similarity_threshold", threshold)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

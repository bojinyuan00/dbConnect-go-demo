package service

import (
	"dbConnect-go-demo/database/models"
	"gorm.io/gorm"
	"sync"
)

type UserService struct {
}

var (
	userService *UserService
	once        sync.Once
)

func NewUserService() *UserService {
	if userService == nil {
		once.Do(func() {
			userService = &UserService{}
		})
	}
	return userService
}

// CreateUser 创建用户
func (u *UserService) CreateUser(db *gorm.DB, name, email, password string) error {
	user := models.User{Name: name, Email: email, Password: password}
	return db.Create(&user).Error
}

// GetUser 查询用户
func (u *UserService) GetUser(db *gorm.DB, id uint) (models.User, error) {
	var user models.User
	err := db.First(&user, id).Error
	return user, err
}

// UpdateUser 更新用户信息
func (u *UserService) UpdateUser(db *gorm.DB, id uint, name, email string) error {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	return db.Save(&user).Error
}

// DeleteUser 删除用户
func (u *UserService) DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&models.User{}, id).Error
}

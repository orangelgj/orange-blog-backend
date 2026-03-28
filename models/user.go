package models

import (
	"gblog/global"
	"time"
)

type User struct {
	ID          uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"unique;not null" json:"username"`
	Password    string    `gorm:"not null" json:"-"` // 密码通常在 JSON 中隐藏
	Email       string    `gorm:"unique;not null" json:"email"`
	Description string    `json:"description"`
	Role        int8      `gorm:"default:2" json:"role"` // 1-管理员, 2-普通用户
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_time"`
}

func GetUserByUsername(username string) (User, error) {
	var user User
	return user, global.DB.Where("username = ?", username).First(&user).Error
}

func GetUserByID(id any) (User, error) {
	var user User
	return user, global.DB.Where("id = ?", id).First(&user).Error
}
func GetUsers() ([]User, error) {
	var users []User
	return users, global.DB.Find(&users).Error
}

// TableName 指定 User 结构体对应的表名
func (User) TableName() string {
	return "user" // 这里返回你实际的表名
}

func CreateUser(m *User) error {
	return global.DB.Create(m).Error
}

func UpdateUsername(id uint32, newUsername string) error {
	return global.DB.Model(&User{}).Where("id = ?", id).Update("username", newUsername).Error
}

func UpdatePassword(id uint32, newPassword string) error {
	return global.DB.Model(&User{}).Where("id = ?", id).Update("password", newPassword).Error
}

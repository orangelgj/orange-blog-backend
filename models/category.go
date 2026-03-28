package models

import (
	"gblog/global"
	"time"
)

type Category struct {
	ID          uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
}

func GetCategories() ([]Category, error) {
	var categories []Category
	err := global.DB.Find(&categories).Error
	return categories, err
}

func (Category) TableName() string {
	return "category"
}

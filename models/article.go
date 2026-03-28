package models

import (
	"gblog/global"
	"time"
)

type Article struct {
	ID         uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"type:longtext;not null" json:"content"`
	Summary    string    `json:"summary"`
	CategoryID uint32    `gorm:"not null;index:idx_category" json:"category_id"`
	Author     string    `gorm:"size:50;not null;index:idx_author" json:"author"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
}

func GetArticleList(pageSize int, pageNum int, categoryID int) ([]Article, error) {
	var articleList []Article

	query := global.DB.Model(&Article{})

	// 如果 categoryID > 0，则添加分类筛选条件
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 分页查询，不计算总数
	err := query.Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Order("create_time DESC").
		Omit("content").
		Find(&articleList).Error

	return articleList, err
}
func (Article) TableName() string {
	return "article"
}
func GetArticleDetail(id uint32) (Article, error) {
	var article Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	return article, err
}

func CreateArticle(article *Article) error {
	return global.DB.Create(article).Error
}

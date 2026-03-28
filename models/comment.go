package models

import (
	"gblog/global"
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ArticleID uint `gorm:"not null;index:idx_article_root;column:article_id" json:"articleId"`
	UserID    uint `gorm:"not null;index:idx_user_id;column:user_id" json:"userId"`

	// 层级与通知核心
	// root_id = 0 表示是顶级评论
	RootId     uint   `gorm:"not null;default:0;index:idx_article_root;column:root_id" json:"rootId"`
	ParentId   uint   `gorm:"not null;default:0;column:parent_id" json:"parentId"`
	ToUserId   uint   `gorm:"not null;default:0;index:idx_to_user;column:to_user_id" json:"toUserId"`
	ToUserName string `gorm:"size:50;column:to_user_name" json:"toUserName"`

	// 内容与状态
	Content   string `gorm:"type:text;not null;column:content" json:"content"`
	Status    int8   `gorm:"type:tinyint(1);default:1;column:status" json:"status"` // 1:正常, 0:待审, -1:删除
	LikeCount uint   `gorm:"not null;default:0;column:like_count" json:"likeCount"`

	// 其他
	IpLocation string    `gorm:"size:64;column:ip_location" json:"ipLocation"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time" json:"createTime"`
}

func (Comment) TableName() string {
	return "comment"
}

func GetRootCommentsByArticleID(articleID uint, pageSize int, pageNum int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	err := global.DB.Model(&Comment{}).
		Where("article_id = ? AND root_id = 0", articleID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = global.DB.Where("article_id = ? AND root_id = 0", articleID).
		Order("create_time DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&comments).Error

	return comments, total, err
}

func GetChildCommentsByRootID(rootID uint, pageSize int, pageNum int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	err := global.DB.Model(&Comment{}).
		Where("root_id = ?", rootID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = global.DB.Where("root_id = ?", rootID).
		Order("create_time ASC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&comments).Error

	return comments, total, err
}

func GetReplyCountByRootID(rootID uint) (int64, error) {
	var count int64
	err := global.DB.Model(&Comment{}).
		Where("root_id = ?", rootID).
		Count(&count).Error
	return count, err
}

func GetFirstReplyByRootID(rootID uint) (Comment, error) {
	var reply Comment
	err := global.DB.Where("root_id = ?", rootID).
		Order("create_time ASC").
		First(&reply).Error
	return reply, err
}

func GetFirstRepliesByRootIDs(rootIDs []uint) (map[uint]Comment, error) {
	var replies []Comment
	replyMap := make(map[uint]Comment)

	if len(rootIDs) == 0 {
		return replyMap, nil
	}

	err := global.DB.Raw(`
		SELECT c1.* FROM comment c1
			INNER JOIN (
				SELECT root_id, MIN(id) as min_id
				FROM comment
				WHERE root_id IN ?
				GROUP BY root_id
			) c2 ON c1.id = c2.min_id
	`, rootIDs).Scan(&replies).Error

	for _, reply := range replies {
		replyMap[reply.RootId] = reply
	}
	return replyMap, err
}

func CreateComment(comment *Comment) error {
	return global.DB.Create(comment).Error
}

func DeleteComment(commentID uint) error {
	return global.DB.Model(&Comment{}).
		Where("id = ?", commentID).
		Update("status", -1).Error
}

func GetCommentByID(commentID uint) (Comment, error) {
	var comment Comment
	err := global.DB.Where("id = ?", commentID).First(&comment).Error
	return comment, err
}

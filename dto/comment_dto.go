package dto

import "time"

// CommentUserDTO 简化的用户信息
type CommentUserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	//Avatar   string `json:"avatar"`
}

// ReplyDTO 评论回复的详细信息
type ReplyDTO struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	UserID     uint      `json:"userId"`
	UserName   string    `json:"userName"`
	ToUserID   uint      `json:"toUserId"`
	ToUserName string    `json:"toUserName"`
	CreateTime time.Time `json:"createTime"`
	LikeCount  uint      `json:"likeCount"`
	IpLocation string    `json:"ipLocation"`
}

// RootCommentDTO 楼层主（根评论）信息
type RootCommentDTO struct {
	ID           uint      `json:"id"`
	Content      string    `json:"content"`
	UserID       uint      `json:"userId"`
	UserName     string    `json:"userName"`
	CreateTime   time.Time `json:"createTime"`
	LikeCount    uint      `json:"likeCount"`
	ReplyCount   int64     `json:"replyCount"`
	IpLocation   string    `json:"ipLocation"`
	PreviewReply *ReplyDTO `json:"previewReply"`
}

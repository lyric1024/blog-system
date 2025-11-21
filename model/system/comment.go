package system

import (
	"time"

	"gorm.io/gorm"
)

// 评论表
type Comment struct {
	ID        uint `gorm:"primarykey" json:"ID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Content string `json:"content"`
	UserID uint `json:"userID"`
	PostID uint `json:"postID"`
}
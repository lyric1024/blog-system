package system

import (
	"time"

	"gorm.io/gorm"
)

// 文章表
type Post struct {
	ID        uint `gorm:"primarykey" json:"ID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title string `json:"title"`
	Content string `json:"content"`
	UserID uint `json:"userID"`
	
}
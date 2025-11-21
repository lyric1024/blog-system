package service

import (
	"github.com/lyric1024/blog-system/model/system"
	"github.com/lyric1024/blog-system/pkg/errors"
	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (commentService *CommentService) CreateComment(userID, postID uint, content string) (*system.Comment, error) {
	comment := &system.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := commentService.db.Create(comment).Error; err != nil {
		return nil, errors.Internal("新增评论失败", err)
	}

	return comment, nil
}

func (commentService *CommentService) ListComment(postID uint) ([]system.Comment, error) {
	comment := []system.Comment{}

	if err := commentService.db.Where("post_id = ? ", postID).Find(&comment).Error; err != nil {
		return nil, errors.Internal("查询评论失败", err)
	}

	return comment, nil
}

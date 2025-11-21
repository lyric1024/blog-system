package service

import (
	"github.com/lyric1024/blog-system/model/system"
	"github.com/lyric1024/blog-system/pkg/errors"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (postService *PostService) CreatePost(userID uint, title, content string) (*system.Post, error) {
	post := &system.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	if err := postService.db.Create(post).Error; err != nil {
		return nil, errors.Internal("创建文章失败", err)
	}

	return post, nil
}

func (postService *PostService) QueryPostList(postID *uint) ([]system.Post, error) {
	posts := []system.Post{}

	db := postService.db
	if postID != nil {
		db = db.Where("id = ?", *postID)
	}
	if err := db.Find(&posts).Error; err != nil {
		return nil, errors.Internal("查询文章失败", err)
	}

	return posts, nil
}

func (postService *PostService) UpdatePost(userID, postID uint, content string) (*system.Post, error) {
	post := &system.Post{}

	// 查询postID
	if err := postService.db.Where("id = ?", postID).First(post).Error; err != nil {
		return nil, errors.Internal("查询文章失败", err)
	}
	// 判断归属
	if post.UserID != userID {
		return nil, errors.Unauthorized("当前用户没有权限", nil)
	}

	if err := postService.db.Model(post).Update("content = ?", content).Error; err != nil {
		return nil, errors.Internal("更新文章失败", err)
	}

	return post, nil
}

func (postService *PostService) DeletePost(userID, postID uint) error {
	post := &system.Post{}

	// 查询postID
	if err := postService.db.Where("id = ?", postID).First(post).Error; err != nil {
		return errors.Internal("查询文章失败", err)
	}
	// 判断归属
	if post.UserID != userID {
		return errors.Unauthorized("当前用户没有权限", nil)
	}

	if err := postService.db.Delete(post).Error; err != nil {
		return errors.Internal("删除文章失败", err)
	}

	return nil
}

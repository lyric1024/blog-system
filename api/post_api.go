package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lyric1024/blog-system/model/common/response"
	"github.com/lyric1024/blog-system/pkg/errors"
	"github.com/lyric1024/blog-system/service"
)

type PostApi struct {
	postService *service.PostService
}

func NewPostApi(postService *service.PostService) *PostApi {
	return &PostApi{postService: postService}
}

func (pApi *PostApi) CreatePostApi(c *gin.Context) {
	userID, exsit := c.Get("userID")
	if !exsit {
		c.Error(errors.Unauthorized("未授权", nil))
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("参数错误", err))
		return
	}

	post, err := pApi.postService.CreatePost(userID.(uint), req.Content, req.Title)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"post": post})

}

func (pApi *PostApi) ListApi(c *gin.Context) {
	_, exsit := c.Get("userID")
	if !exsit {
		c.Error(errors.Unauthorized("未授权", nil))
		return
	}

	var req struct {
		PostID *uint `json:"postID"`
	}

	//
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("参数错误", err))
		return
	}

	posts, err := pApi.postService.QueryPostList(req.PostID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"posts": posts})
}

func (pApi *PostApi) UpdateApi(c *gin.Context) {
	userID, exsit := c.Get("userID")
	if !exsit {
		c.Error(errors.Unauthorized("未授权", nil))
		return
	}

	var req struct {
		PostID  uint   `json:"postID" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("参数错误", err))
		return
	}

	post, err := pApi.postService.UpdatePost(userID.(uint), req.PostID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"post": post})
}

func (pApi *PostApi) DeleteApi(c *gin.Context) {
	userID, exsit := c.Get("userID")
	if !exsit {
		c.Error(errors.Unauthorized("未授权", nil))
		return
	}

	var req struct {
		PostID uint `json:"postID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("参数错误", err))
		return
	}

	err := pApi.postService.DeletePost(userID.(uint), req.PostID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"postID": req.PostID})
}

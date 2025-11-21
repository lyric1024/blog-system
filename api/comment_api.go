package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lyric1024/blog-system/model/common/response"
	"github.com/lyric1024/blog-system/pkg/errors"
	"github.com/lyric1024/blog-system/service"
)

type CommentApi struct {
	commentService *service.CommentService
}

func NewCommentApi(commentService *service.CommentService) *CommentApi {
	return &CommentApi{commentService: commentService}
}

// 新增评论
func (cApi *CommentApi) CreateCommentApi(c *gin.Context) {
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

	comment, err := cApi.commentService.CreateComment(userID.(uint), req.PostID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"comment": comment})
}

// 查询评论
func (cApi *CommentApi) ListApi(c *gin.Context) {
	_, exsit := c.Get("userID")
	if !exsit {
		c.Error(errors.Unauthorized("未授权", nil))
		return
	}

	var req struct {
		PostID uint `json:"postID" binding:"required"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.Error(errors.BadRequest("参数错误", err))
		return
	}

	comments, err := cApi.commentService.ListComment(req.PostID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"comments": comments})
}

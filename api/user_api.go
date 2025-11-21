package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lyric1024/blog-system/model/common/response"
	"github.com/lyric1024/blog-system/model/system"
	"github.com/lyric1024/blog-system/pkg/errors"
	"github.com/lyric1024/blog-system/pkg/jwt"
	"github.com/lyric1024/blog-system/pkg/logger"
	"github.com/lyric1024/blog-system/service"
	"go.uber.org/zap"
)

type UserAPI struct {
	userService *service.UserService
}

func NewUserAPI(userService *service.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// 注册
func (uApi *UserAPI) RegistApi(c *gin.Context) {
	var user system.User

	err := c.ShouldBind(&user)
	if err != nil {
		c.Error(errors.BadRequest("注册函数参数错误", err))
		return
	}

	// 检查用户是否存在
	if err := uApi.userService.ValidUser(user.UserName, user.Email); err != nil {
		c.Error(err)
		return
	}
	// 创建用户
	saveUser, err := uApi.userService.CreateUser(user)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"user": saveUser})
}

// 登录
func (uApi *UserAPI) LoginApi(c *gin.Context) {

	var req struct {
		UserName string `json:"userName" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(errors.BadRequest("登录函数参数错误", err))
		return
	}

	user, err := uApi.userService.Login(req.UserName, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	// 生成jwt token
	token, err := jwt.GetToken(user.ID)
	if err != nil {
		c.Error(errors.Internal("token生成失败", err))
		logger.Error("token生成失败", zap.Error(err), zap.String("userName", req.UserName))
		return
	}
	jwt.SetTokenCookie(c, token)

	response.Success(c, gin.H{"user": user, "token": token})
}

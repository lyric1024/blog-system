package router

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/lyric1024/blog-system/api"
	"github.com/lyric1024/blog-system/model/common/response"
	"github.com/lyric1024/blog-system/pkg/errors"
	"github.com/lyric1024/blog-system/pkg/jwt"
	"github.com/lyric1024/blog-system/pkg/logger"
	"github.com/lyric1024/blog-system/service"
)

func InitApiRouter(r *gin.Engine, db *gorm.DB) {

	// 健康检查接口
	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// 初始化用户service和api
	userService := service.NewUserService(db)
	userApi := api.NewUserAPI(userService)
	// 无须认证
	public := r.Group("/api")
	{
		public.POST("/regist", userApi.RegistApi)
		public.POST("/login", userApi.LoginApi)
	}

	// 初始化文章service和api
	postService := service.NewPostService(db)
	postApi := api.NewPostApi(postService)
	// 初始化评论service和api
	commentService := service.NewCommentService(db)
	commentApi := api.NewCommentApi(commentService)

	// 需认证
	privated := r.Group("/api").Use(JWTMiddleware())
	{
		privated.POST("/post/create", postApi.CreatePostApi)
		privated.POST("/post/list", postApi.ListApi)
		privated.POST("/post/update", postApi.UpdateApi)
		privated.POST("/post/delete", postApi.DeleteApi)
	}
	{
		privated.POST("/comment/create", commentApi.CreateCommentApi)
		privated.POST("/comment/list", commentApi.ListApi)
	}
}

// jwt中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		// 1
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
		// 2
		if tokenString == "" {
			if cookieToken, err := c.Cookie("token"); err != nil && cookieToken != "" {
				tokenString = cookieToken
			}
		}
		// 3
		if tokenString == "" {
			queryToken := c.Query("token")
			if queryToken != "" {
				tokenString = queryToken
			}
		}

		if tokenString == "" {
			c.Error(errors.Unauthorized("缺少有效的认证凭证（请提供 Authorization 头、token Cookie 或 token 查询参数）", nil))
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.Error(errors.Unauthorized("无效或过期的 token", nil))
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)

		//
		c.Next()
	}
}

// 全局错误
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 执行后

		// 封装err
		for _, err := range c.Errors {
			ErrAndRespond(c, err)
			return
		}
	}
}
func ErrAndRespond(c *gin.Context, err error) {
	code := int(errors.ErrCodeInternal)
	msg := "服务器内部错误"

	if e, ok := err.(*errors.AppError); ok {
		code = int(e.Code)
		msg = e.Message
		logger.Warn("AppError", zap.Int("code", code), zap.String("message", msg))
	} else { // 非AppError
		logger.Error("【ERROR】not AppError",
			zap.String("path", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)
	}

	response.Fail(c, code, msg)
}

// 日志中间件zap
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next() // 执行后

		// 记录日志
		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 根据状态码决定日志级别
		if statusCode >= 500 {
			logger.Error("HTTP 5xx", fields...)
		} else if statusCode >= 400 {
			logger.Warn("HTTP 4xx", fields...)
		} else {
			logger.Info("HTTP access", fields...)
		}

	}
}

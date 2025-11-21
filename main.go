package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/lyric1024/blog-system/configs"
	"github.com/lyric1024/blog-system/model/system"
	"github.com/lyric1024/blog-system/pkg/jwt"
	"github.com/lyric1024/blog-system/pkg/logger"
	"github.com/lyric1024/blog-system/router"
	"github.com/spf13/viper"
)

func main() {

	initBlogSystem()

	fmt.Println("blog system exec end ...")
}

func initBlogSystem() {
	// ======== 1. 加载config配置文件  ========
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 当前目录
	// 读取config.yaml
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("配置文件config.yaml未找到, 请在项目根目录配置")
		} else {
			log.Fatalf("读取文件config.yaml失败: %+v", err)
		}
	}

	// 解析到结构体
	var config configs.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("解析配置失败: %+v", err)
	}
	// 初始化日志
	logger.Init(config.Log.Level, config.Log.OutputFile)
	defer logger.Sync()

	// ======== 2. 初始化db  ========
	dsn := config.Mysql.Dsn()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("初始化db失败: %+v", err)
		os.Exit(1)
	}
	Run(db)

	// ======== 3. 初始化路由  ========
	jwt.Init(config.Jwt.Secret, config.Jwt.ExpireTime) // jwt全局单例

	r := gin.Default()
	r.Use(router.ErrorHandle(), router.RequestLogger()) // 全局错误处理 日志中间件zap
	router.InitApiRouter(r, db)

	// ======== 4. 启动服务  ========
	log.Printf("http server端口号为 %s", config.System.Port)
	server := &http.Server{
		Addr:    config.System.Port,
		Handler: r,
	}

	// ======== 5. 在goroutine启动服务  ========
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http服务启动失败: %+v", err)
		}
	}()

	// ======== 6. 实现优雅关闭  ========
	closeFuncs := []func(){
		// 1. 关闭http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("http server关闭失败: %+v \n", err)
			} else {
				log.Println("http server已关闭")
			}
		},
		// 2. 关闭DB
		func() {
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("从gorm中获取sqlDB失败: %+v \n", err)
			}

			if err := sqlDB.Close(); err != nil {
				log.Printf("mysql连接关闭失败: %+v \n", err)
			} else {
				log.Println("mysql连接已关闭")
			}
		},
	}

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("收到系统关闭信号，启动关闭逻辑... \n")

	for _, closeFn := range closeFuncs {
		closeFn()
	}
	log.Printf("系统已成功关闭... \n")
}

// 创建初始表
func Run(db *gorm.DB) {
	db.AutoMigrate(&system.User{}, &system.Post{}, &system.Comment{})
}

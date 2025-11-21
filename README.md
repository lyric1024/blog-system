# blog-system
博客系统后端（gin & gorm）

## 概要
这是一个基于 Gin + GORM 的简单博客系统后端。包含用户注册、登录（JWT）、文章和评论的增删改查接口。

## 运行环境
- Go 1.18+（推荐使用 1.20+）
- MySQL（示例使用 `test` 数据库）

## 仓库结构（关键目录）
- `main.go` - 启动程序与初始化逻辑
- `router/` - 路由与中间件（定义了公开与受保护的 API）
- `api/` - 各资源的 HTTP 处理器（user/post/comment）
- `service/` - 业务逻辑层
- `model/` - GORM 模型
- `config.yaml` - 示例配置（项目根目录）
- `tests/` - 测试脚本

## 配置（config.yaml）
项目根目录已包含 `config.yaml`，示例内容：

```yaml
# system configuration
system:
	port: ":8080"
	
# jwt
jwt:
    secret: "blogSystem"
    expiretime: 168 #  24 * 7

# mysql配置
mysql:
	host: "127.0.0.1"
	port: "3306"
	config: "charset=utf8mb4&parseTime=True&loc=Local"
	db-name: "test"
	username: "username"
	password: "pwd"

log:
	level: "debug"
	output-file: ""  # 留空表示输出到stdout
```

请根据本地环境修改 `mysql` 配置（用户名、密码、数据库名等）。

## 依赖安装 & 编译
在项目根目录运行：

```bash
go mod download
go build -o blog-system ./
```

也可以直接用 `go run` 运行（适合开发）：

```bash
go run main.go
```

注意：程序会在启动时根据 `config.yaml` 初始化数据库连接并调用 `AutoMigrate`，会尝试创建 `users/posts/comments` 等表。

## HTTP 接口 & 测试（curl）
路由在 `router/router.go` 中定义，主要接口如下：

1) 公开接口（无需 JWT）
- POST /api/regist — 用户注册
- POST /api/login  — 用户登录（返回 token）

2) 受保护接口（需携带 token，支持 Authorization: Bearer <token> / token cookie / token query）
- POST /api/post/create  — 新建文章
- POST /api/post/list    — 查询文章（可按 postID 查询单篇）
- POST /api/post/update  — 更新文章
- POST /api/post/delete  — 删除文章
- POST /api/comment/create — 新增评论
- POST /api/comment/list   — 查询评论（按 postID）

## 如何运行并收集真实测试结果
1. 确保 MySQL 可用，并在 `config.yaml` 中配置正确的 `username/password/db-name`。
2. 如果数据库中尚无 `test` 数据库，先创建：

```sql
CREATE DATABASE test CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

3. 启动服务：

```bash
go run main.go
```

4. 在另一终端执行下述脚本来测试请求

```bash
bash tests/run_full_test.sh
```

注意事项与常见问题
- 如果遇到 `token` 相关错误，请确认登录确实返回了 `token` 并在请求中以 `Authorization: Bearer <token>` 形式传入。
- 如果启动时报 DB 连接错误，检查 `config.yaml` 的 MySQL 字段并确保 MySQL 可用。

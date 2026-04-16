## 中文 | [English](readme.md)

- [项目简介](#项目简介)
- [许可证](#许可证)
- [版本记录](#版本记录)
- [安装说明](#安装说明)
  - [克隆项目](#克隆项目)
  - [初始化Go环境与依赖](#初始化Go环境与依赖)
    - [方式一](#方式一)
    - [方式二](#方式二)
  - [启动](#启动)
    - [使用air热更新](#使用air热更新)
  - [编译](#编译)
    - [编译项目](#编译项目)
    - [编译命令行](#编译命令行)
- [目录结构](#目录结构)
- [使用方法](#使用方法)
  - [启动服务](#启动服务)
    - [air热更新](#air热更新)
  - [配置文件](#配置文件)
    - [项目配置](#项目配置)
    - [热更新配置](#热更新配置)
  - [中间件](#中间件)
    - [中间件创建帮助](#中间件创建帮助)
    - [中间件创建](#中间件创建)
    - [限流中间件](#限流中间件)
  - [路由](#路由)
    - [路由创建帮助](#路由创建帮助)
    - [路由创建](#路由创建)
    - [路由列表](#路由列表)
  - [控制器](#控制器)
    - [控制器创建帮助](#控制器创建帮助)
    - [控制器创建](#控制器创建)
  - [模型](#模型)
    - [模型创建帮助](#模型创建帮助)
    - [模型创建](#模型创建)
    - [ORM动态筛选](#ORM动态筛选)
    - [OR条件查询](#OR条件查询)
    - [AND条件查询](#AND条件查询)
    - [json字段查询](#json字段查询)
    - [复杂查询](#复杂查询)
  - [表单验证](#表单验证)
    - [验证创建帮助](#验证创建帮助)
    - [验证创建](#验证创建)
    - [验证规则](#验证规则)
    - [验证场景](#验证场景)
    - [提示信息](#提示信息)
    - [字段翻译](#字段翻译)
    - [自定义验证](#自定义验证)
      - [全局规则](#全局规则)
      - [局部规则](#局部规则)
      - [临时规则](#临时规则)
      - [验证使用](#验证使用)
      - [在控制器中使用](#在控制器中使用)
  - [服务](#服务)
    - [服务创建帮助](#服务创建帮助)
    - [服务创建](#服务创建)
  - [命令行](#命令行)
    - [获取版本](#获取版本)
    - [命令帮助](#命令帮助)
    - [命令列表](#命令列表)
    - [命令创建帮助](#命令创建帮助)
    - [命令创建](#命令创建)
    - [命令结构](#命令结构)
    - [命令注册](#命令注册)
    - [帮助选项](#帮助选项)
    - [执行命令](#执行命令)
    - [编译执行](#编译执行)
  - [缓存](#缓存)
    - [全局缓存](#全局缓存)
    - [Redis缓存](#Redis缓存)
    - [内存缓存](#内存缓存)
    - [磁盘缓存](#磁盘缓存)
  - [事件](#事件)
    - [事件创建帮助](#事件创建帮助)
    - [事件创建](#事件创建)
  - [监听](#监听)
    - [监听创建帮助](#监听创建帮助)
    - [监听创建](#监听创建)
  - [队列](#队列)
    - [队列创建帮助](#队列创建帮助)
    - [队列创建](#队列创建)
    - [队列使用](#队列使用)
    - [消费者列表](#消费者列表)
    - [生产者列表](#生产者列表)
  - [发布事件](#发布事件)
    - [测试事件](#测试事件)
  - [事件列表](#事件列表)
    - [事件监听列表](#事件监听列表)
  - [响应](#响应)
    - [成功响应](#成功响应)
      - [成功提示](#成功提示)
      - [成功数据](#成功数据)
    - [失败响应](#失败响应)
      - [失败错误码](#失败错误码)
      - [失败提示](#失败提示)
      - [失败数据](#失败数据)
  - [日志](#日志)
    - [记录日志](#记录日志)
    - [错误调试](#错误调试)
  - [多语言](#多语言)
    - [目录配置](#目录配置)
    - [常规翻译](#常规翻译)
    - [模版翻译](#模版翻译)
    - [添加语言](#添加语言)
  - [容器服务](#容器服务)
    - [容器服务创建](#容器服务创建)
  - [门面](#门面)
    - [门面创建](#门面创建)
    - [门面使用](#门面使用)
  - [数据库](#数据库)
    - [数据库配置](#数据库配置)
    - [数据库连接](#数据库连接)
    - [数据库搜索](#数据库搜索)
  - [swagger文档](#swagger文档)

# 项目简介
> 基于`Golang`语言框架`Go Gin`开发的轻量级框架, 开箱即用, 设计灵感基于`Laravel`、`ThinPHP`等主流`PHP`框架, 项目架构目录层次分明, 初学者的福音, 框架默认集成了`门面`、`容器服务`、`jwt`、`日志`、`中间件`、`缓存`、`验证器`、`事件`、`路由`、`队列(kafka、rabbitmq)`、`redis`、`命令行`等,支持多语言,开发简单易于上手, 方便扩展。
## 项目地址
- Github: https://github.com/dsxwk/gin-admin.git
- Gitee: https://gitee.com/dsxwk/gin-admin.git

## Gin框架介绍
> Gin是一个用Go语言编写的Web框架。它具有简单、快速、高效等特点, 被广泛应用于Go语言的Web开发中。

## Gin框架的特性
- 快速: Gin框架基于标准库net/http，使用goroutine和channel实现异步处理，提高性能。
- 简单: Gin框架提供了一系列的API和中间件，使得开发人员可以快速构建Web应用程序。
- 高效: Gin框架使用sync.Pool来缓存对象，减少内存分配和释放，提高性能。
> Golang Gin 是一个轻量级且高效的 Golang Web 框架。它具有高性能、易用性和灵活性等特点，被广泛应用于各种 Web 应用程序的开发。

# 许可证
- 📘 开源版: 遵循 AGPL-3.0，仅供学习、研究及非商业用途。
- 💼 商业版: 如需闭源或商业使用，请联系作者📧  [25076778@qq.com] 获取商业授权。

# 版本记录
> - 最新版本: v2.0.10
> - [版本更新详细记录](version_history_zh.md)

# 安装说明
> 项目基于Golang 1.25.2版本开发, 低版本可能存在版本差异, 建议版本 >= 1.25.2。
## 克隆项目
```bash
$ git clone https://github.com/dsxwk/gin.git
$ cd gin
$ mv dev.config.yaml.example dev.config.yaml
```
## 初始化Go环境与依赖
### 方式一
```bash
$ go env -w GOPROXY=https://goproxy.cn,direct
$ go generate ./...
```
### 方式二
```bash
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
# $ go get -u
$ go mod tidy
# $ go mod download
$ go mod vendor
```
## 启动
```bash
$ go run main.go
```
### 使用air热更新
```bash
$ go install github.com/air-verse/air@latest
$ air
```

## 编译
### 编译项目
```bash
$ go build main.go
$ ./main
```

### 编译命令行
```bash
$ go build ./cmd/cli.go
$ ./cli demo:command --args=11

执行命令: demo:command, 参数: 11
```

# 目录结构
```
├── app                                 # 应用程序
│   ├── command                         # 命令
│   ├── controller                      # 控制器
│   ├── event                           # 事件
│   ├── facade                          # 门面
│   ├── listener                        # 监听
│   ├── middleware                      # 中间件
│   ├── model                           # 模型
│   ├── provider                        # 容器服务
│   ├── queue                           # 消息队列
│   ├──├── kafka                        # Kafka
│   ├──├──├── consumer                  # 消费者
│   ├──├──├── producer                  # 生产者
│   ├──├── rabbitmq                     # Rabbitmq
│   ├──├──├── consumer                  # 消费者
│   ├──├──├── producer                  # 生产者
│   ├── request                         # 验证器
│   ├── service                         # 服务
├── bootstrap                           # 初始化文件 
├── cmd                                 # 命令行工具
│   ├── cli.go                          # 命令行工具入口文件
├── common                              # 公共模块
│   ├── base                            # 基类
│   ├── ctxkey                          # 上下文键名
│   ├── errcode                         # 错误码
│   ├── flag                            # 特殊符合
│   ├── response                        # 响应
│   ├── template                        # 模版
├── config                              # 配置文件
├── database                            # 数据库测试文件
├── docs                                # 文档
├── pkg                                 # 工具包
│   ├──├── cache                        # 缓存
│   ├──├── cli                          # 命令行
│   ├──├── debugger                     # 调试
│   ├──├── eventbus                     # 事件
│   ├──├── foundation                   # 服务提供者
│   ├──├── http                         # http请求
│   ├──├── lang                         # 多语言
│   ├──├── logger                       # 日志
│   ├──├── message                      # 消息事件
│   ├──├── orm                          # orm工具
│   ├──├── queue                        # 队列
│   ├──├── time                         # 时间处理
├── public                              # 静态资源
├── router                              # 路由
├── storage                             # 存储
│   ├── cache                           # 磁盘缓存
│   ├── logs                            # 日志
│   ├── locales                         # 翻译文件
│   ├──├── en                           # 英文翻译
│   ├──├── zh                           # 中文翻译
├── tests                               # 测试用例
├── tmp                                 # 临时文件
├── vendor                              # 依赖包
├── .air.linux.toml                     # air配置文件
├── .air.toml                           # air配置文件
├── .gitignore                          # git忽略文件
├── config.yaml                         # 默认配置文件
├── dev.config.yaml                     # 本地环境配置文件
├── go.mod                              # go mod
├── LICENSE                             # 开源协议
├── main.go                             # 入口文件
├── readme.md                           # 英文文档
├── readme_zh.md                        # 中文文档
├── VersionHistoryEn.md                 # 版本记录英文文档
└── VersionHistoryZn.md                 # 版本记录中文文档
```

# 使用方法
## 启动服务
```bash
$ go run main.go
```
### air热更新
```bash
$ go install github.com/air-verse/air@latest
$ air

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ v1.62.0, built with Go go1.24.2

watching .
watching app
watching app\command
watching app\controller
...
...
[GIN-debug] GET    /api/v1/user/:id          --> gin/app/controller/v1.(*UserController).Detail-fm (6 handlers)
应用:                  gin
环境:                  dev
端口:                  8080
数据库:                gin
🌐 Local Address:      http://127.0.0.1:8080
🌐 Network Address:    http://192.168.8.54:8080
👉 Local Swagger:      http://127.0.0.1:8080/swagger/index.html
👉 Network Swagger:    http://192.168.8.54:8080/swagger/index.html
👉 Local Test API:     http://127.0.0.1:8080/ping
👉 Network Test API:   http://192.168.8.54:8080/ping
 SUCCESS  Gin server started successfully!
```

## 配置文件
### 项目配置
> `config.yaml`为默认配置文件, 可自行修改。`dev.config.yaml`对应本地环境配置, 通过以下app.env文件配置环境变量来切换环境
> ```
> app:
>   env: dev # dev|testing|production dev=本地环境 testing=测试环境 production=生产环境
> ```

### 热更新配置
> `.air.toml`为Windows环境下默认配置文件, `.air.linux.toml`为Linux环境下默认配置文件。可自行根据项目整体需要自行修改。

## 中间件
> `middleware`目录下为中间件目录, 可自行添加中间件, 并在`router/root.go`文件中注册中间件。
### 中间件创建帮助
```bash
$ go run ./cmd/cli.go make:middleware -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:middleware  中间件创建

Options:
  -f, --file  文件路径, 如: auth    required:true
  -d, --desc  描述, 如: 权限中间件  required:false
```

### 中间件创建
```bash
$ go run ./cmd/cli.go make:middleware --file=auth --desc=授权中间件
```

### 限流中间件
> `middleware/rate_limit.go`文件中定义了全局限流中间件, 支持全局用户接口限流、ip接口限流以及全局限流。
```go
package router

import (
    "gin/app/middleware"
    "github.com/gin-gonic/gin"
)

var rateLimitMiddleware middleware.RateLimit

// LoadRouters 加载路由
func LoadRouters(router *gin.Engine) {
    // 全局限流
    group := router.Group("", rateLimitMiddleware.Handle())
	r := group.Group("") 
	{
        r.GET("/global-test1", func(c *gin.Context) {
            err := errcode.NewError(0, "global test1")
            response.Success(c, &err)
        })

        r.GET("/global-test2", func(c *gin.Context) {
            err := errcode.NewError(0, "global test2")
            response.Success(c, &err)
        })
    }

	// 指定接口限流
    // 用户限流
    // r 每秒产生多少token
    // burst 桶容量
    userGroup := router.Group("", rateLimitMiddleware.UserRateLimit(1, 1))
	r1 := userGroup.Group("")
	{
		r1.GET("/user-test1", func(c *gin.Context) {
            err := errcode.NewError(0, "user test1")
            response.Success(c, &err)
        })

        r1.GET("/user-test2", func(c *gin.Context) {
            err := errcode.NewError(0, "user test2")
            response.Success(c, &err)
        })
    }

    // 指定接口限流
    // ip限流
    // r 每秒产生多少token
    // burst 桶容量
    ipGroup := router.Group("", rateLimitMiddleware.IpRateLimit(1, 1))
	r2 := ipGroup.Group("")
	{
		r2.GET("/ip-test1", func(c *gin.Context) {
            err := errcode.NewError(0, "ip test1")
            response.Success(c, &err)
        })

        r2.GET("/ip-test2", func(c *gin.Context) {
            err := errcode.NewError(0, "ip test2")
            response.Success(c, &err)
		})
    }
}
```

## 路由
> `router/root.go` 文件中定义了全局路由规则可自行修改,  一般情况只需要默认即可。
### 路由创建帮助
```bash
$ go run ./cmd/cli.go make:router -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:router  路由创建

Options:
  -f, --file  文件路径, 如: user      required:true
  -d, --desc  路由描述, 如: 用户路由   required:false
```

### 路由创建
```bash
$ go run ./cmd/cli.go make:router --file=user --desc=用户路由
```
```go
package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// UserRouter 用户路由
type UserRouter struct{}

func init() {
	Register(&UserRouter{})
}

// RegisterRoutes 注册路由
func (r *UserRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		user v1.UserController
	)

	router := routerGroup.Group("api/v1")
	{
		// 列表
		router.GET("/user", user.List)
		// 创建
		router.POST("/user", user.Create)
		// 更新
		router.PUT("/user/:id", user.Update)
		// 删除
		router.DELETE("/user/:id", user.Delete)
		// 详情
		router.GET("/user/:id", user.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *UserRouter) IsAuth() bool {
	return true
}

```

### 路由列表
```bash
$ go run ./cmd/cli.go route:list

---------------------------------------------------------
Method   Path                                Handler
---------------------------------------------------------
POST     /api/v1/login                       gin/app/controller/v1.(*LoginController).Login
GET      /api/v1/user                        gin/app/controller/v1.(*UserController).List
POST     /api/v1/user                        gin/app/controller/v1.(*UserController).Create
GET      /api/v1/user/:id                    gin/app/controller/v1.(*UserController).Detail
PUT      /api/v1/user/:id                    gin/app/controller/v1.(*UserController).Update
DELETE   /api/v1/user/:id                    gin/app/controller/v1.(*UserController).Delete
GET      /ping                               gin/router.LoadRouters
GET      /public/*filepath                   github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler
HEAD     /public/*filepath                   github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler
GET      /swagger/*any                       github.com/swaggo/gin-swagger.CustomWrapHandler
---------------------------------------------------------
总计 10 条路由
```

## 控制器
### 控制器创建帮助
```bash
$ go run ./cmd/cli.go make:controller -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:controller  控制器创建

Options:
  -f, --file      文件路径, 如: v1/user  required:true
  -F, --function  方法名称, 如: list     required:false
  -m, --method    请求方式, 如: get      required:false
```

### 控制器创建
```bash
$ go run ./cmd/cli.go make:controller --file=v1/test --router=/test --method=get --desc=列表 --function=list
```
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

// List 列表
// @Router /test [get]
func (s *TestController) List(c *gin.Context) {
    // Define your function here
    s.Success(c, errcode.Success().WithMsg("Test Msg").WithData([]string{}))
}
```

## 模型
### 模型创建帮助
```bash
$ go run ./cmd/cli.go make:model -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:model  模型创建

Options:
  -t, --table       表名, 如: user 或 user,menu  required:true
  -p, --path        输出目录, 如: api/user       required:false
  -c, --camel       json字段是否使用驼峰          required:false
  -C, --connection  数据库连接                   required:false
```

### 模型创建
> 支持创建同时多个模型文件, 如需创建多个模型文件, 传人的表名参数请使用逗号分隔, 如: user,menu
```bash
$ go run ./cmd/cli.go make:model --table='user,menu' --path=api/user --camel=true --connection=mysql
# go run ./cmd/cli.go make:model-old --table=user --path=api/user --camel=true --connection=sqlsrv
```
```go
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package user

import "gin/app/model"

const TableNameUser = "user"

// User 用户表
type User struct {
	ID        int64            `gorm:"column:id;type:int(10) unsigned;primaryKey;autoIncrement:true;comment:ID" json:"id"`           // ID
	Avatar    string           `gorm:"column:avatar;type:varchar(255);not null;comment:头像" json:"avatar"`                            // 头像
	Username  string           `gorm:"column:username;type:varchar(10);not null;comment:用户名" json:"username"`                        // 用户名
	FullName  string           `gorm:"column:full_name;type:varchar(20);not null;comment:姓名" json:"fullName"`                        // 姓名
	Email     string           `gorm:"column:email;type:varchar(50);not null;comment:邮箱" json:"email"`                               // 邮箱
	Password  string           `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"password"`                        // 密码
	Nickname  string           `gorm:"column:nickname;type:varchar(50);not null;comment:昵称" json:"nickname"`                         // 昵称
	Gender    int64            `gorm:"column:gender;type:tinyint(1) unsigned;not null;comment:性别 1=男 2=女" json:"gender"`             // 性别 1=男 2=女
	Age       int64            `gorm:"column:age;type:int(10) unsigned;not null;comment:年龄" json:"age"`                              // 年龄
	Status    int64            `gorm:"column:status;type:tinyint(3) unsigned;not null;default:1;comment:状态 1=启用 2=停用" json:"status"` // 状态 1=启用 2=停用
	CreatedAt *model.DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt"`                                // 创建时间
	UpdatedAt *model.DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`                                // 更新时间
	DeletedAt *model.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" swaggerignore:"true"`                                // 删除时间
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

// Connection 数据库连接名称
func (*User) Connection() string {
  return "mysql"
}
```

## ORM动态筛选
> 通过`post`或者`get`传递`query`|`body`参数`__search`根据列表字段动态指定查询条件,`__search`类型为`map[string]interface{}` 如: __search={"and":[{"username":"test"},{"age":18}]}, __search={"or":[{"username":"test"},{"age":18}]}. 支持or、and、in、not in、between、not between、like、left like、right like、is not null、is null、gt、gte、lt、lte、exist、not exist、json_contains、json_extract等条件,不区分大小写.参数支持两种模式{"username":"admin"}或者{"username":["like", "admin"]},字段名为mysql where条件的关键字时自动根据条件构建sql语句.
### OR条件查询
```http
GET /api/v1/user?__search={"or":[{"username":"test"},{"age":18}]} // {"or":[{"username":["=", "test"]},{"age":["=", 18]}]}
```
```sql
SELECT * FROM `user` WHERE (username = 'test' OR age = 18)
```

### AND条件查询
```http
GET /api/v1/user?__search={"and":[{"username":"test"},{"age":18}]} // {"and":[{"username":["=", "test"]},{"age":["=", 18]}]}
```
```sql
SELECT * FROM `user` WHERE (username = 'test' AND age = 18)
```

### json字段查询
```http
GET /api/v1/menu?__search={"or":[{"and":[{"createdAt":[">","2025-01-01"]},{"createdAt":["<","2026-01-01"]},{"name":""},{"$.meta.icon":["=","ele-Collection"]}]}]}
```
```sql
 SELECT * FROM `menu` WHERE ((((menu.created_at > '2025-01-01') AND (menu.created_at < '2026-01-01') AND (menu.name = '') AND (JSON_EXTRACT(meta, '$.icon') = 'ele-Collection'))))
```

### 复杂查询
```http
GET /api/v1/user?__search={"or":[{"and":[{"createdAt":[">","2025-01-01"]},{"createdAt":["<","2026-01-01"]},{"not exist":{"userRoles.name":"admin"}}]},{"username":"admin"}]}
```
```sql
 SELECT * FROM `user` WHERE ((((user.created_at > '2025-01-01') AND (user.created_at < '2026-01-01') AND (NOT EXISTS (SELECT 1 FROM user_roles WHERE user_roles.user_id = user.id AND user_roles.name = 'admin'))) OR (user.username = 'admin')))
```

## 表单验证
### 验证创建帮助
```bash
$ go run ./cmd/cli.go make:request -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:request  验证请求创建

Options:
  -f, --file        文件路径, 如: user  required:true
  -t, --table       表名, 如: roles    required:false
  -d, --desc        描述               required:false
  -c, --camel       字段是否使用驼峰     required:false
  -C, --connection  数据库连接          required:false
```

### 验证创建
```bash
$ go run ./cmd/cli.go make:request --file=roles --table=roles --desc=用户请求验证
```
```go
package request

import (
  "errors"
  "gin/common/base"
  "github.com/gookit/validate"
)

// Roles 角色请求验证
type Roles struct {
  base.BaseRequest
  ID     int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
  Name   string `json:"name" form:"name" validate:"required|max:255" label:"角色名称"`
  Desc   string `json:"desc" form:"desc" validate:"required|max:255" label:"角色描述"`
  Status int64  `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
  PageListValidate
}

// Validate 请求验证
func (s Roles) Validate(data Roles, scene string) error {
  v := validate.Struct(data, scene)
  if !v.Validate(scene) {
    return errors.New(v.Errors.One())
  }
  return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Roles) ConfigValidation(v *validate.Validation) {
  scenes := validate.SValues{
    "list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
    "create": []string{"Name", "Desc", "Status"},
    "update": []string{"ID", "Name", "Desc", "Status"},
    "detail": []string{"ID"},
    "delete": []string{"ID"},
  }
  v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Roles) Messages() map[string]string {
  return validate.MS{
    "required":    "字段 {field} 必填",
    "int":         "字段 {field} 必须为整数",
    "gt":          "字段 {field} 必须大于 0",
    "max":         "字段 {field} 长度不能超过 255",
    "Page.gt":     "页码必须大于 0",
    "PageSize.gt": "每页数量必须大于 0",
  }
}

// Translates 字段翻译
func (s Roles) Translates() map[string]string {
  return validate.MS{
    "ID":       "ID",
    "Name":     "角色名称",
    "Desc":     "角色描述",
    "Status":   "状态 1=启用 2=停用",
    "Page":     "页码",
    "PageSize": "每页数量",
  }
}
```

### 验证规则
> 更多规则请查看 [gookit/validate](https://github.com/gookit/validate)
```go
package request

// Roles 角色请求验证
type Roles struct {
  base.BaseRequest
  ID     int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
  Name   string `json:"name" form:"name" validate:"required|max:255" label:"角色名称"`
  Desc   string `json:"desc" form:"desc" validate:"required|max:255" label:"角色描述"`
  Status int64  `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
  PageListValidate
}
```

### 验证场景
```go
package request

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Roles) ConfigValidation(v *validate.Validation) {
  scenes := validate.SValues{
    "list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
    "create": []string{"Name", "Desc", "Status"},
    "update": []string{"ID", "Name", "Desc", "Status"},
    "detail": []string{"ID"},
    "delete": []string{"ID"},
  }
  v.WithScenes(scenes)
}
```

### 提示信息
```go
package request

// Messages 验证器错误消息
func (s Roles) Messages() map[string]string {
	return validate.MS{
        "required":                     "字段 {field} 必填",
        "int":                          "字段 {field} 必须为整数",
        "PageListValidate.Page.gt":     "字段 {field} 需大于 0",
        "PageListValidate.PageSize.gt": "字段 {field} 需大于 0",
	}
}
```

### 字段翻译
```go
package request

// Translates 字段翻译
func (s Roles) Translates() map[string]string {
  return validate.MS{
    "ID":       "ID",
    "Name":     "角色名称",
    "Desc":     "角色描述",
    "Status":   "状态 1=启用 2=停用",
    "Page":     "页码",
    "PageSize": "每页数量",
  }
}
```

### 自定义验证
#### 全局规则
> 全局规则只需要在入口文件`main.go`中定义, 适用于所有验证器, 无需重复定义。
```go
package main

import (
	"github.com/gookit/validate"
)

// 初始化时注册
func init() {
	validate.AddValidator("is_even", func(val any, rule string) bool {
		num, ok := val.(int)
		if !ok {
			return false
		}
		return num%2 == 0
	})
}
```

#### 局部规则
```go
package request

// ValidateIsEven 定义局部规则方法(命名规则：Validate<规则名>)
func (s User) ValidateIsEven(val any) bool {
	num := val.(int)
	return num%2 == 0
}
```

#### 临时规则
```go
package request

// Validate 请求验证
func (s User) Validate(data User, scene string) error {
	v := validate.Struct(data, scene)
	v.AddValidator("is_even", func(val any, rule string) bool {
        num, ok := val.(int)
        if !ok {
            return false
        }
        return num%2 == 0
    })
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}
```

#### 验证使用
```go
package request

type User struct {
    Age int `json:"gender" validate:"required|is_even" label:"年龄"`
}
```

#### 在控制器中使用
```go
package v1

import (
  "gin/app/facade"
  "gin/app/model"
  "gin/app/request"
  "gin/app/service"
  "gin/common/base"
  "gin/common/errcode"
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/copier"
  "strconv"
)

type UserController struct {
  base.BaseController
  service service.UserService
}

// List 列表
// @Tags 用户管理
// @Summary 列表
// @Description 用户列表
// @Param token header string true "认证Token"
// @Param page query string true "页码"
// @Param pageSize query string true "分页大小"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.User}} "登录成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/user [get]
func (s *UserController) List(c *gin.Context) {
  var (
    ctx = c.Request.Context()
    req request.User
  )

  s.service.WithContext(ctx)

  err = c.ShouldBind(&req)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(err.Error()))
    return
  }

  // 验证
  err = req.Validate(req, "List")
  if err != nil {
    s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
    return
  }

  res, err := s.service.List(req)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(facade.Lang.T(ctx, err.Error(), nil)))
    return
  }

  s.Success(c, errcode.Success().WithData(res))
}
```

## 服务
### 服务创建帮助
```bash
$ go run ./cmd/cli.go make:service -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:service  服务创建

Options:
  -f, --file      文件路径, 如: v1/user  required:true
  -F, --function  方法名称, 如: list     required:false
  -d, --desc      描述, 如: 列表         required:false
exit status 3
```

### 服务创建
```bash
$ go run ./cmd/cli.go make:service -f=user --function=list --desc="列表"
```

## 命令行
### 获取版本
```bash
$ go run ./cmd/cli.go --version # -v
Gin Cli v2.0.0
```

### 命令帮助
```bash
$ go run ./cmd/cli.go -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Available commands:
consumer:
  consumer:list    消费者列表
db:
  db:migrate       数据迁移
  db:rollback      数据回滚
  db:seed          数据填充
demo:
  demo:command     test-demo
event:
  event:list       事件列表
listener:
  listener:list    事件监听列表
make:
  make:command     命令创建
  make:controller  控制器创建
  make:event       创建事件
  make:facade      创建门面
  make:listener    创建监听
  make:middleware  中间件创建
  make:migration   生成数据库迁移模板
  make:model       模型创建
  make:model-old   模型创建old
  make:provider    创建服务提供者
  make:queue       创建消息队列(Kafka/RabbitMQ)
  make:request     验证请求创建
  make:router      路由创建
  make:seed        生成数据库seeder模板
  make:service     服务创建
producer:
  producer:list    生产者列表
route:
  route:list       路由列表

Options:
  -f, --format     The output format (txt, json) [default: txt]
  -h, --help       Display help for the given command
  -v, --version    Display CLI version
```

### 命令列表
```bash
$ go run ./cmd/cli.go --format=json # -f=json
{
  "commands": [
    {
      "description": "消费者列表",
      "name": "consumer:list"
    },
    {
      "description": "数据迁移",
      "name": "db:migrate"
    },
    {
      "description": "数据回滚",
      "name": "db:rollback"
    },
    {
      "description": "数据填充",
      "name": "db:seed"
    },
    {
      "description": "test-demo",
      "name": "demo:command"
    },
    {
      "description": "事件列表",
      "name": "event:list"
    },
    {
      "description": "事件监听列表",
      "name": "listener:list"
    },
    {
      "description": "命令创建",
      "name": "make:command"
    },
    {
      "description": "控制器创建",
      "name": "make:controller"
    },
    {
      "description": "创建事件",
      "name": "make:event"
    },
    {
      "description": "创建门面",
      "name": "make:facade"
    },
    {
      "description": "创建监听",
      "name": "make:listener"
    },
    {
      "description": "中间件创建",
      "name": "make:middleware"
    },
    {
      "description": "生成数据库迁移模板",
      "name": "make:migration"
    },
    {
      "description": "模型创建",
      "name": "make:model"
    },
    {
      "description": "模型创建old",
      "name": "make:model-old"
    },
    {
      "description": "创建服务提供者",
      "name": "make:provider"
    },
    {
      "description": "创建消息队列(Kafka/RabbitMQ)",
      "name": "make:queue"
    },
    {
      "description": "验证请求创建",
      "name": "make:request"
    },
    {
      "description": "路由创建",
      "name": "make:router"
    },
    {
      "description": "生成数据库seeder模板",
      "name": "make:seed"
    },
    {
      "description": "服务创建",
      "name": "make:service"
    },
    {
      "description": "生产者列表",
      "name": "producer:list"
    },
    {
      "description": "路由列表",
      "name": "route:list"
    }
  ],
  "version": "Gin Cli v2.0.0"
}
```

## 命令创建帮助
```bash
$ go run ./cmd/cli.go make:command -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:command  命令创建

Options:
  -f, --file  文件路径, 如: cronjob/demo  required:true
  -n, --name  命令名称, 如: demo-test     required:false
  -d, --desc  描述, 如: command-desc      required:false
```

## 命令创建
```bash
$ go run ./cmd/cli.go make:command --file=cronjob/demo --name=demo-test --desc=command-desc
```

## 命令结构
> 生成命令后，应为`Name()` 和 `Description()` 方法定义适当的值。当在显示命令列表时，将使用这些属性。 `Name()` 方法还允许你定义命令的输入期望值。 `Execute()` 执行命令时将调用该方法。你可以将命令逻辑放在此方法中。 让我们看一个示例命令。
```go
package cronjob

import (
	"gin/common/base"
	"gin/pkg/cli"
	"github.com/fatih/color"
)

type DemoCommand struct {
	base.BaseCommand
}

func (m *DemoCommand) Name() string {
    return "demo-test"
}

func (m *DemoCommand) Description() string {
	return "command-desc"
}

func (m *DemoCommand) Help() []base.CommandOption {
	return []base.CommandOption{
        {
            base.Flag{
                Short: "a",
                Long:  "args",
            },
            "示例参数, 如: arg1",
            true,
        },
    }
}

func (m *DemoCommand) Execute(args []string) {
    values := m.ParseFlags(m.Name(), args, m.Help())
    color.Green("执行命令: %s %s", m.Name(), m.FormatArgs(values))
}

func init() {
	cli.Register(&DemoCommand{})
}

```

## 命令注册
> `./cmd/cli.go` 默认注册了 `gin/app/command` 目录下的 `command` 包的所有命令，如果你注册的命令不是一个包，可以在 `./cmd/imports/import.go` 中添加导入包的路径。
```go
package main

import (
    _ "gin/cmd/imports"
    "gin/config"
    "gin/pkg/cli"
)

func main() {
    _ = config.NewConfig()
	cli.Execute()
}
```

## 帮助选项
> 命令选项参数使用 `base.CommandOption` 结构体来定义。 `base.CommandOption` 结构体包含两个属性： `Flag` 和 `Description`。 `Flag` 属性用于定义命令选项的标志，可以是短标志（如 `-a`）或长标志（如 `--args`）。 `Description` 属性用于定义命令选项的描述。 `base.CommandOption` 结构体还包含一个 `Required` 属性，用于指定命令选项是否为必需的。同时该方法支持控制台 `--help` 参数，自动生成帮助信息。
```go
func (m *DemoCommand) Help() []base.CommandOption {
	return []base.CommandOption{
        {
            base.Flag{
                Short: "a",
                Long:  "args",
            },
            "示例参数, 如: arg1",
            true,
        },
    }
}
```
```bash
$ go run ./cmd/cli.go demo-test -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  demo-test  command-desc

Options:
  -a, --args  示例参数, 如: arg1  required:true
```

## 执行命令
```bash
$ go run ./cmd/cli.go demo:command --args=arg1
 SUCCESS  执行命令: demo:command --args=arg1
```

## 编译执行
```bash
$ go build ./cmd/cli.go
$ ./cli demo:command --args=arg1
```

# 缓存
> 默认使用 `memory` 作为缓存驱动, 支持自定义扩展。默认支持`内存缓存`、`Redis缓存`、`磁盘缓存`三种模式, 可使用全局缓存也可单独使用任意缓存。全局缓存默认只集成了`Set`、`Get`、`Delete`、`Expire`公共方法如需使用更多可以单独使用,你也可以自己集成。
## 全局缓存
> 配置全局缓存可通过`yaml`配置文件中的`cache.driver`配置进行切换，也可以动态切换。
```go
package controller

import (
	"fmt"
    "gin/app/facade"
    "gin/common/base"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test() {
    // Set 设置缓存	
    key := "test_key"
    value := "test_value"
	cache := facade.Cache.Store()
	cache = facade.Cache.Store("redis")
    err := cache.Set(key, value, time.Second*10)
	if err != nil {
	    // 处理错误	
    }
	
    // Get 获取缓存
    key := "test_key"
    value := "test_value"
    result, ok := cache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete 删除缓存
	key := "test_key"
	err := cache.Delete(key)
	if err != nil {
        // 处理错误	
    }
	
	// Expire 获取缓存过期时间
	key := "test_key"
    val, expireAt, ok, err := cache.Expire(key)
	if err != nil {
	    // 处理错误
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
}
```

## Redis缓存
```go
package controller

import (
	"fmt"
    "gin/app/facade"
    "gin/common/base"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test() {
    // Set 设置缓存	
    key := "test_key"
    value := "test_value"
	redisCache := facade.Cache.Store("redis")
    err := redisCache.Set(key, value, time.Second*10)
	if err != nil {
	    // 处理错误	
    }
	
    // Get 获取缓存
    key := "test_key"
    value := "test_value"
    result, ok := redisCache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete 删除缓存
	key := "test_key"
	err := redisCache.Delete(key)
	if err != nil {
        // 处理错误	
    }
	
	// Expire 获取缓存过期时间
	key := "test_key"
    val, expireAt, ok, err := redisCache.Expire(key)
	if err != nil {
	    // 处理错误
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
	
	// ... 其他
}
```

## 内存缓存
```go
package controller

import (
	"fmt"
    "gin/app/facade"
    "gin/common/base"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test() {
    // Set 设置缓存	
    key := "test_key"
    value := "test_value"
    memoryCache := facade.Cache.Store("memory")
    err := memoryCache.Set(key, value, time.Second*10)
	if err != nil {
	    // 处理错误	
    }
	
    // Get 获取缓存
    key := "test_key"
    value := "test_value"
    result, ok := memoryCache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete 删除缓存
	key := "test_key"
	err := memoryCache.Delete(key)
	if err != nil {
        // 处理错误	
    }
	
	// Expire 获取缓存过期时间
	key := "test_key"
    val, expireAt, ok, err := memoryCache.Expire(key)
	if err != nil {
	    // 处理错误
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
	
	// ... 其他
}
```

## 磁盘缓存
```go
package controller

import (
    "fmt"
    "gin/app/facade"
    "gin/common/base"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test() {
    // Set 设置缓存	
    key := "test_key"
    value := "test_value"
    diskCache := facade.Cache.Store("disk")
    err := diskCache.Set(key, value, time.Second*10)
    if err != nil {
        // 处理错误	
    }
    
    // Get 获取缓存
    key := "test_key"
    value := "test_value"
    result, ok := diskCache.Get(key)
    if ok {
        println(result) // test_value	
    }
    
    // Delete 删除缓存
    key := "test_key"
    err := diskCache.Delete(key)
    if err != nil {
        // 处理错误	
    }
    
    // Expire 获取缓存过期时间
    key := "test_key"
    val, expireAt, ok, err := diskCache.Expire(key)
    if err != nil {
        // 处理错误
    }
    if ok {
        fmt.Println(val) // test_value
        fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
    
    // ... 其他	
}
```

# 事件
## 事件创建帮助
```bash
$ go run ./cmd/cli.go make:event -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:event  创建事件

Options:
  -f, --file  文件路径, 如: login/test  required:true
  -n, --name  事件名称, 如: test-event  required:false
  -d, --desc  事件描述, 如: 测试事件     required:false
```

## 事件创建
```bash
$ go run ./cmd/cli.go make:event -f=user_login -n='user.login' -d=用户登录事件
```
```go
package event

// UserLoginEvent 事件数据
type UserLoginEvent struct {
	UserId   int64
	Username string
}

// Name 事件名称
func (u UserLoginEvent) Name() string {
	return "user.login"
}

// Description 事件描述
func (u UserLoginEvent) Description() string {
	return "用户登录事件"
}

```

# 监听
## 监听创建帮助
```bash
$ go run ./cmd/cli.go make:listener -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:listener  创建监听

Options:
  -f, --file   文件路径, 如: login/test  required:true
  -e, --event  事件数据, 如: UserLogin   required:true
```

## 监听创建
```bash
$ go run ./cmd/cli.go make:listener -f=user_login -e=UserLoginEvent
```
```go
package listener

import (
    "fmt"
    "gin/app/event"
    "gin/common/base"
    "gin/pkg/eventbus"
    "github.com/goccy/go-json"
    "time"
)

type UserLoginListener struct{}

func (l *UserLoginListener) Handle(e base.Event) {
    ev, ok := e.(event.UserLoginEvent)
    if !ok {
        return
    }
  
	data, _ := json.Marshal(e)
	fmt.Printf("收到事件: %s 事件描述: %s 事件数据: %s, 时间: %s\n", ev.Name(), ev.Description(), data, time.Now().Format("2006-01-02 15:04:05"))
}

func init() {
	eventbus.Register(&UserLoginListener{}, event.UserLoginEvent{})
}

```

# 队列
> 执行队列创建命令会根据队列类型同时创建消费者和生产者, 如: kafka会创建kafka消费者和生产者, rabbitmq会创建rabbitmq消费者和生产者. 你只需要完善消费者当中`Handle`方法完善你的业务逻辑即可, 支持自动错误重试以及延迟队列.
## 队列创建帮助
```bash
$ go run ./cmd/cli.go make:queue -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:queue  消息队列创建

Options:
  -t, --type      队列类型, 如: kafka或rabbitmq  required:true
  -n, --name      队列文件名称, 如: order_create  required:true
  -d, --isDelay   是否延迟队列, 如: true或false   required:false
  -T, --topic     队列主题, 如: kafka_demo       required:false
  -k, --key       消息键, 如: kafka_demo         required:false
  -g, --group     消费组, 如: kafka_demo         required:false
  -q, --queue     队列名, 如: rabbitmq_demo      required:false
  -e, --exchange  交换机, 如: rabbitmq_demo      required:false
  -r, --routing   路由键, 如: rabbitmq_demo      required:false
  -R, --retry     错误重试次数, 如: 3             required:false
  -m, --delayMs   延迟毫秒, 如: 10000            required:false
```

## 队列创建
```bash
$ go run ./cmd/cli.go make:queue --type=rabbitmq --name=rabbitmq_demo --queue=rabbitmq_demo --exchange=rabbitmq_demo --routing=rabbitmq_demo 
```
```go
package consumer

import (
  "gin/app/facade"
  "gin/common/base"
  "gin/config"
  "gin/pkg"
  "gin/pkg/logger"
  "gin/pkg/queue"
)

// RabbitmqDemoConsumer RabbitMQ普通消费者
type RabbitmqDemoConsumer struct {
  *base.RabbitmqConsumer
}

// NewRabbitmqDemoConsumer 创建消费者实例
func NewRabbitmqDemoConsumer() *RabbitmqDemoConsumer {
  cfg := facade.Config.Get()
  log := facade.Log.Logger()
  bus := facade.Message.GetBus()

  // 创建RabbitMQ连接
  mq, err := base.NewRabbitMQ(cfg, log, bus)
  if err != nil {
    log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
    return nil
  }

  return &RabbitmqDemoConsumer{
    RabbitmqConsumer: &base.RabbitmqConsumer{
      Mq:           mq,
      Queue:        "rabbitmq_demo",
      Exchange:     "rabbitmq_demo_exchange",
      Routing:      "rabbitmq_demo",
      IsDelayQueue: false,
      Retry:        3,
    },
  }
}

// Name 消费者名称
func (c *RabbitmqDemoConsumer) Name() string {
  return "rabbitmq_demo"
}

// Start 启动消费者
func (c *RabbitmqDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
  c.RabbitmqConsumer.Start(c)
  log.Info(pkg.Sprintf("RabbitMQ消费者启动成功: %s", c.Name()))
  return nil
}

// Stop 停止消费者
func (c *RabbitmqDemoConsumer) Stop() error {
  return c.RabbitmqConsumer.Stop()
}

// Enabled 是否启用
func (c *RabbitmqDemoConsumer) Enabled(cfg *config.Config) bool {
  return cfg.Rabbitmq.Enabled
}

// Status 消费者状态
func (c *RabbitmqDemoConsumer) Status() queue.ConsumerStatus {
  return c.RabbitmqConsumer.Status()
}

// Handle 处理消息的业务逻辑
func (c *RabbitmqDemoConsumer) Handle(msg string) error {
  facade.Log.Info(pkg.Sprintf("RabbitMq Received Msg: %s", msg))
  // todo 处理业务逻辑
  return nil
}

// init 注册消费者到注册表
func init() {
  queue.GetConsumerRegistry().Register(NewRabbitmqDemoConsumer())
}
```
```go
package producer

import (
  "gin/app/facade"
  "gin/common/base"
  "gin/pkg"
  "gin/pkg/queue"
)

// RabbitmqDemoProducer RabbitMQ普通生产者
type RabbitmqDemoProducer struct {
  *base.RabbitmqProducer
}

// NewRabbitmqDemoProducer 创建生产者实例
func NewRabbitmqDemoProducer() *RabbitmqDemoProducer {
  cfg := facade.Config.Get()
  log := facade.Log.Logger()
  bus := facade.Message.GetBus()

  mq, err := base.NewRabbitMQ(cfg, log, bus)
  if err != nil {
    log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
    return nil
  }

  return &RabbitmqDemoProducer{
    RabbitmqProducer: &base.RabbitmqProducer{
      Mq:           mq,
      Queue:        "rabbitmq_demo",
      Exchange:     "rabbitmq_demo_exchange",
      Routing:      "rabbitmq_demo",
      IsDelayQueue: false,
    },
  }
}

func (p *RabbitmqDemoProducer) Name() string {
  return "rabbitmq_demo"
}

func init() {
  queue.GetProducerRegistry().Register(NewRabbitmqDemoProducer())
}
```

## 队列使用
> 消费者启动项目时自动注册在容器中无限额外启动,生产者直接使用门面初始化即可使用。
```go
package controller

import (
    "fmt"
    "gin/app/facade"
    "gin/common/base"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test() {
    // 获取生产者
    producer := facade.Queue.Producer("rabbitmq_demo")
	_ = producer.Publish(ctx, []byte(`{"orderId":111, "message":"message 111"}`))
}
```

## 消费者列表
```bash
$ go run ./cmd/cli.go consumer:list

┌────────────────────────────────────────────────────────────┐
│ 消费者名称             描述                                  │
├────────────────────────────────────────────────────────────┤
│ kafka_delay_demo       kakfa延迟队列消费者                   │
│ kafka_demo             kakfa普通队列消费者                   │
│ rabbitmq_delay_demo    rabbitmq延迟队列消费者                │
│ rabbitmq_demo          rabbitmq普通队列消费者                │
└────────────────────────────────────────────────────────────┘
总计 4 个消费者
```

## 生产者列表
```bash
$ go run ./cmd/cli.go producer:list

┌────────────────────────────────────────────────────────────┐
│ 生产者名称               描述                                │
├────────────────────────────────────────────────────────────┤
│ kafka_delay_demo       kafka延迟队列生产者                   │
│ kafka_demo             kafka普通队列生产者                   │
│ rabbitmq_delay_demo    rabbitmq延迟队列生产者                │
│ rabbitmq_demo          rabbitmq普通队列生产者                │
└────────────────────────────────────────────────────────────┘
总计 4 个生产者
```

# 发布事件
```go
package v1

import (
  "gin/app/event"
  "gin/app/facade"
  "gin/app/model"
  "gin/app/request"
  "gin/app/service"
  "gin/common/base"
  "gin/common/errcode"
  "github.com/gin-gonic/gin"
)

type LoginController struct {
	base.BaseController
    service service.LoginService
}

// Token token信息
type Token struct {
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	TokenExpire        int64  `json:"tokenExpire" example:"7200"`
	RefreshTokenExpire int64  `json:"refreshTokenExpire" example:"172800"`
}

type LoginResponse struct {
	Token Token `json:"token"`
	User  model.User
}

// Login 登录
// @Tags 登录相关
// @Summary 登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param data body request.UserLogin true "登录参数"
// @Success 200 {object} errcode.SuccessResponse{data=LoginResponse} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/login [post]
func (s *LoginController) Login(c *gin.Context) {
  var (
    ctx = c.Request.Context()
    req request.Login
  )

  s.service.WithContext(ctx)

  err := c.ShouldBind(&req)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(err.Error()))
    return
  }

  // 验证
  err = req.Validate(req, "Login")
  if err != nil {
    s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
    return
  }

  userModel, err := s.service.Login(req.Username, req.Password)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(facade.Lang.T(ctx, err.Error(), nil)))
    return
  }

  err, userModel, accessToken, refreshToken, tokenExpire, refreshTokenExpire := s.service.Login(req.Username, req.Password)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(facade.Lang.T(ctx, err.Error(), nil)))
    return
  }

  // 发布事件
  facade.Event.Publish(ctx, event.UserLoginEvent{
    UserId:   userModel.ID,
    Username: userModel.Username,
  })

  s.Success(
    c, errcode.Success().WithMsg(
      facade.Lang.T(ctx, "login.success", map[string]interface{}{
        "name": userModel.Username,
      }),
    ).WithData(LoginResponse{
      Token{
        AccessToken:        accessToken,
        RefreshToken:       refreshToken,
        TokenExpire:        tokenExpire,
        RefreshTokenExpire: refreshTokenExpire,
      },
      userModel,
    }),
  )
}
```

## 测试事件
```bash
$ POST /api/v1/login HTTP/1.1
Host: 127.0.0.1:8080
Accept-Language: en-Us
Content-Type: application/json
Content-Length: 56

{
    "username": "admin",
    "password": "123456"
}

收到事件: user.login 事件描述: 用户登录事件 事件数据: {"UserId":1,"Username":"admin"}, 时间: 2025-11-04 15:32:12
```

# 事件列表
```bash
$ go run ./cmd/cli.go event:list

┌────────────────────────────────────────────────────────────┐
│ 事件名称               描述                                  │
├────────────────────────────────────────────────────────────┤
│ user.login            用户登录事件                           │
└────────────────────────────────────────────────────────────┘
总计 1 个事件
```

## 事件监听列表
```bash
$ go run ./cmd/cli.go listener:list

┌────────────────────────────────────────────────────────────┐
│ 事件名称               描述                                  │
├────────────────────────────────────────────────────────────┤
│ user.login             用户登录事件                          │
│                      ├─ *listener.TestListener             │
│                      └─ *listener.UserLoginListener        │
└────────────────────────────────────────────────────────────┘
总计 1 个事件 2 个监听
```

# 响应
## 成功响应
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Success(c, errcode.Success())
}
```

### 成功提示
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Success(c, errcode.Success().WithMsg("Success"))
}
```

### 成功数据
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Success(c, errcode.Success().WithData([]string{"test data"}))
}
```

## 失败响应
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Error(c, errcode.SystemError())
}
```

### 失败错误码
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Error(c, errcode.SystemError().WithCode(500))
}
```

### 失败提示
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Error(c, errcode.SystemError().WithMsg("System Error"))
}
```

### 失败数据
```go
package v1

import (
    "gin/common/base"
    "gin/common/errcode"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    return s.Error(c, errcode.SystemError().WithData([]string{"test data"}))
}
```

# 日志
> 使用 `zap` 包实现日志记录，日志文件存放路径为 `storage/logs`, 默认日志级别为 `debug`, 返回错误码不为0时自动记录日志TraceId、堆栈、sql、http、redis等调用信息, 也可以直接调用日志记录也会自动记录调试信息。配置文件`yaml`中`log.access`支持是否自动记录请求日志，如若开启会自动记录请求日志。
```json
{
  "level": "info",
  "timestamp": "2025-11-17 16:35:09.402",
  "caller": "middleware/logger.go:83",
  "msg": "Access Log",
  "traceId": "fa505122-d31e-4d4f-a05c-13c1641d6c6c",
  "ip": "127.0.0.1",
  "path": "/api/v1/login",
  "method": "POST",
  "params": {
    "password": "1234561",
    "username": "admin"
  },
  "ms": 59,
  "debugger": {
    "Sql": [
      {
        "ms": 2.5008,
        "rows": 1,
        "sql": "SELECT * FROM `user` WHERE username = 'admin' AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1"
      }
    ],
    "Cache": [],
    "Http": [],
    "Mq": [],
    "ListenerEvent": []
  }
}
```

## 记录日志
> 已封装在门面中, 日志级别支持debug、info、warn、error、dPanic、panic、fatal, 默认为`debug`。
```go
package v1

import (
    "gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    facade.Log.Error("System Error")
}
```

## 错误调试
> 使用公共返回错误以及调用WithDebugger()方法时会自动记录日志TraceId、堆栈、sql、http、redis等调用信息, 可根据debug调试信息和trace堆栈信息调试, 日志文件存放路径为 `storage/logs`。
```go
package v1

import (
	"gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
  ctx := c.Request.Context()
  facade.Log.WithDebugger(ctx).Error("System Error")
}
```
```json
{
  "level": "error",
  "timestamp": "2025-11-17 16:35:09.401",
  "caller": "response/response.go:60",
  "msg": "Login Password Error",
  "traceId": "fa505122-d31e-4d4f-a05c-13c1641d6c6c",
  "ip": "127.0.0.1",
  "path": "/api/v1/login",
  "method": "POST",
  "params": {
    "password": "1234561",
    "username": "admin"
  },
  "ms": 58,
  "debugger": {
    "Sql": [
      {
        "ms": 2.5008,
        "rows": 1,
        "sql": "SELECT * FROM `user` WHERE username = 'admin' AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1"
      }
    ],
    "Cache": [],
    "Http": [],
    "Mq": [],
    "ListenerEvent": []
  },
  "stackTrace": "gin/common/response.Error\n\tE:/www/dsx/www-go/gin/common/response/response.go:60\ngin/common/base.(*BaseController).Error\n\tE:/www/dsx/www-go/gin/common/base/base_controller.go:25\ngin/app/controller/v1.(*LoginController).Login\n\tE:/www/dsx/www-go/gin/app/controller/v1/login.go:67\ngithub.com/gin-gonic/gin.(*Context).Next\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/context.go:192\ngin/router.init.Cors.Handle.func2\n\tE:/www/dsx/www-go/gin/app/middleware/cors.go:30\ngithub.com/gin-gonic/gin.(*Context).Next\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/context.go:192\ngin/router.init.Logger.Handle.func1\n\tE:/www/dsx/www-go/gin/app/middleware/logger.go:76\ngithub.com/gin-gonic/gin.(*Context).Next\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/context.go:192\ngithub.com/gin-gonic/gin.CustomRecoveryWithWriter.func1\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/recovery.go:92\ngithub.com/gin-gonic/gin.(*Context).Next\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/context.go:192\ngithub.com/gin-gonic/gin.LoggerWithConfig.func1\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/logger.go:249\ngithub.com/gin-gonic/gin.(*Context).Next\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/context.go:192\ngithub.com/gin-gonic/gin.(*Engine).handleHTTPRequest\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/gin.go:689\ngithub.com/gin-gonic/gin.(*Engine).ServeHTTP\n\tE:/www/dsx/www-go/gin/vendor/github.com/gin-gonic/gin/gin.go:643\nnet/http.serverHandler.ServeHTTP\n\tE:/go-sdk/go1.25.2/src/net/http/server.go:3340\nnet/http.(*conn).serve\n\tE:/go-sdk/go1.25.2/src/net/http/server.go:2109"
}
```

# 多语言
> 多语言已集成在门面和容器中,语言支持 `zh` 和 `en` 两种语言, 可支持自定义扩展。语言传输默认在`header`中传输 `Accept-Language` 参数, 如 `zh` 或 `en`, 不区分大小写, 不传递默认语言为 `zh`。
## 目录配置
> 翻译文件存放路径为 `storage/locales`, 默认语言为 `zh`, 多个语言用逗号分隔。语言存放在对应的语言目录下不区分子目录, 如中文就放在`storage/locales/zh`下,可以支持任意目录下的`json`和`yaml`格式文件。
```yaml
# 翻译配置
i18n:
  dir: "storage/locales" # 翻译文件存放路径
  lang: "zh,en" # 默认语言,多个语言用逗号分隔
```

## 常规翻译
```go
package controller

import (
    "fmt"
    "gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context)  {
	ctx := c.Request.Context()
    trans := facade.Lang.T(ctx, "login.username", nil)
	fmt.Println(trans) // 输出: 用户名, 英文输出: Username
}
```

## 模版翻译
> 翻译文件中支持模版翻译, 如 `{{.name}}`, 使用 `map[string]interface{}` 传递参数。
```json
[
  {
    "id": "login.success",
    "translation": "{{.name}},登录成功"
  }
]
```
```go
package controller

import (
	"fmt"
    "gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context)  {
    ctx := c.Request.Context()
    trans := facade.Lang.T(ctx, "login.success", map[string]interface{}{
        "name": "admin",
    }),
	fmt.Println(trans) // 输出: admin,登录成功 英文输出: admin,Login Success
}
```

## 添加语言
> 在 `storage/locales` 目录下添加对应语言目录, 如 `en`, 然后在目录下添加翻译文件, 翻译文件支持 `json` 和 `yaml` 格式, 翻译文件中 `id` 为唯一标识, `translation` 为翻译内容, 翻译文件中可以添加任意数量的翻译内容。配置语言支持需调整配置文件i18n.lang参数。
```yaml
# 翻译配置
i18n:
  dir: "storage/locales" # 翻译文件存放路径
  lang: "zh,en" # 默认语言,多个语言用逗号分隔
```

# 容器服务
> 容器会在启动时自动加载注册，关闭时自动释放。
## 容器服务创建
> 同模型、控制器等使用命令行创建,具体参考之前文档。

# 门面
## 门面创建
> 同模型、控制器等使用命令行创建,具体参考之前文档。

## 门面使用
> 项目以默认集成了日志、数据库、缓存、限流等门面，当前以缓存为示例。数据库、缓存、http请求、队列绑定了上下文会记录到调试日志中。
```go
package controller

import (
    "gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context)  {
    ctx := c.Request.Context()
    cache := facade.Cache.Store()
	redisCache := facade.Cache.Store('redis')   // 或 facade.Cache.Redis()
	// 绑定上下文
    redisCache = redisCache.WithContext(ctx)
	memoryCache := facade.Cache.Store('memory') // 或 facade.Cache.Memory()
	diskCache := facade.Cache.Store('disk')     // 或 facade.Cache.Disk()
    // 其他门面使用 ...
}
```

# 数据库
> 数据库通过容器服务初始化,默认集成了mysql、pgsql、sqlite、sqlsrv,可配置默认数据库以及通过Connection方法指定数据库连接。
## 数据库配置
```yaml
# 数据库
databases:
  db-connection: mysql # 默认数据库
  # 慢查询的时间(ms) 超过这个时间会记录到日志中
  slow-query-duration: 3000ms # 3秒(time.Duration)

# Mysql数据库
mysql:
  driver: mysql
  # host: "username:password@tcp(127.0.0.1:3306)/databaseName?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: gin
  # 慢查询的时间(ms) 超过这个时间会记录到日志中
  slow-query-duration: 3000ms # 3秒(time.Duration)

# Postgresql数据库
pgsql:
  driver: pgsql
  host: 127.0.0.1
  port: 5432
  username: testuser
  password: 123456
  database: testdb
  # 慢查询的时间(ms) 超过这个时间会记录到日志中
  slow-query-duration: 3000ms # 3秒(time.Duration)

# sqlite数据库
sqlite:
  driver: sqlite
  path: storage/data/gin.db
  # 慢查询的时间(ms) 超过这个时间会记录到日志中
  slow-query-duration: 3000ms # 3秒(time.Duration)

# sqlsrv数据库
sqlsrv:
  driver: sqlsrv
  host: 127.0.0.1
  port: 1433
  username: root
  password: root
  database: gin
  # 慢查询的时间(ms) 超过这个时间会记录到日志中
  slow-query-duration: 3000ms # 3秒(time.Duration)
```

## 数据库连接
> 使用上下文非必须，如果不绑定上下文则日志不会记录sql记录。
```go
package controller

import (
    "gin/app/facade"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context)  {
    ctx := c.Request.Context()
    // 默认连接
    db := facade.DB.Connection()
    // 使用上下文
    db1 := facade.DB.Connection().WithContext(ctx)
    // 连接pgsql
    db2 := facade.DB.Connection("pgsql").WithContext(ctx)
    // 连接sqlsrv
    db3 := facade.DB.Connection("sqlsrv").WithContext(ctx)
    // todo ...
}
```

## 数据库搜索
> 配合文档中的ORM动态筛选示例使用。
```go
package controller

import (
    "gin/app/facade"
    "gin/app/model"
    "gin/app/request"
    "gin/common/base"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    var (
        ctx = c.Request.Context()
		req request.User
        m   []model.User
		db  = facade.DB.Connection().WithContext(ctx)
	)

    err := c.ShouldBind(&req)
    if err != nil {
        s.Error(c, errcode.SystemError().WithMsg(err.Error()))
        return
    }

    db = db.Model(&model.User{})
	
    if req.Search != nil {
        whereSql, args, err := orm.BuildCondition(req.Search, db, model.User{})
        if err != nil {
            s.Error(c, errcode.SystemError().WithMsg(err.Error()))
            return
        }
    
        if whereSql != "" {
            db = db.Where(whereSql, args...)
        }
    }

    err = db.Offset(10).Limit(10).Order("id DESC").Find(&m).Error
	if err != nil {
		s.Error(c, errcode.SystemError().WithMsg(err.Error()))
		return
    }
	
	s.Success(c, m)
}

```

# swagger文档
```bash
$ go install github.com/swaggo/swag/cmd/swag@latest
$ swag init -g main.go # --exclude cli,app/service
2025/10/23 16:26:42 Generate swagger docs....
2025/10/23 16:26:42 Generate general API Info, search dir:./
2025/10/23 16:26:43 Generating request.UserLogin
2025/10/23 16:26:43 Generating errcode.SuccessResponse
2025/10/23 16:26:43 Generating v1.LoginResponse
2025/10/23 16:26:43 Generating v1.Token
2025/10/23 16:26:43 Generating model.User
2025/10/23 16:26:43 Generating model.DateTime
2025/10/23 16:26:43 Generating errcode.ArgsErrorResponse
2025/10/23 16:26:43 Generating errcode.SystemErrorResponse
2025/10/23 16:26:43 Generating request.PageData
2025/10/23 16:26:43 Generating request.UserCreate
2025/10/23 16:26:43 Generating request.UserUpdate
2025/10/23 16:26:43 Generating request.UserDetail
2025/10/23 16:26:43 create docs.go at docs/docs.go
2025/10/23 16:26:43 create swagger.json at docs/swagger.json
2025/10/23 16:26:43 create swagger.yaml at docs/swagger.yaml
```

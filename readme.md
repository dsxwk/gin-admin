## English | [中文](readme_zh.md)

- [Project Introduction](#Project-Introduction)
- [License](#License)
- [Version History](#Version-History)
- [Installation Instructions](#Installation-Instructions)
  - [Clone Project](#Clone-Project)
  - [Initialize Go Environment And Dependencies](#Initialize-Go-Environment-And-Dependencies)
    - [Method One](#Method-One)
    - [Method Two](#Method-Two)
  - [Start](#Start)
    - [Use Air Hot Update](#Use-Air-Hot-Update)
  - [Compile](#Compile)
    - [Compile Project](#Compile-Project)
    - [Compile Command](#Compile-Command)
- [Directory Structure](#Directory-Structure)
- [Instructions For Use](#Instructions-For-Use)
  - [Start Service](#Start-Service)
    - [Air Hot Update](#Air-Hot-Update)
  - [Configuration File](#Configuration-File)
    - [Project Configuration](#Project-Configuration)
    - [Hot Update Configuration](#Hot-Update-Configuration)
  - [Middleware](#Middleware)
    - [Middleware Creation Help](#Middleware-Creation-Help)
    - [Middleware Creation](#Middleware-Creation)
    - [Rate Limit Middleware](#Rate-Limit-Middleware)
  - [Route](#Route)
    - [Route Creation Help](#Route-Creation-Help)
    - [Route Creation](#Route-Creation)
    - [Route List](#Route-List)
  - [Controller](#Controller)
    - [Controller Creation Help](#Controller-Creation-Help)
    - [Controller Creation](#Controller-Creation)
  - [Model](#Model)
    - [Model Creation Help](#Model-Creation-Help)
    - [Model Creation](#Model-Creation)
    - [ORM Dynamic Filtering](#ORM-Dynamic-Filtering)
    - [OR Condition Query](#OR-Condition-Query)
    - [AND Condition Query](#AND-Condition-Query)
    - [JSON Field Query](#JSON-Field-Query)
    - [Complex Condition Query](#Complex-Condition-Query)
  - [Form Validation](#Form-Validation)
    - [Validator Creation Help](#Validator-Creation-Help)
    - [Validator Creation](#Validator-Creation)
    - [Validator Rules](#Validator-Rules)
    - [Validator Scenes](#Validator-Scenes)
    - [Prompt Message](#Prompt-Message)
    - [Field Translation](#Field-Translation)
    - [Custom Validation](#Custom-Validation)
      - [Global Rules](#Global-Rules)
      - [Local Rules](#Local-Rules)
      - [Temporary Rules](#Temporary-Rules)
      - [Validator Usage](#Validator-Usage)
      - [Used In The Controller](#Used-In-The-Controller)
  - [Service](#Service)
    - [Service Creation Help](#Service-Creation-Help)
    - [Service Creation](#Service-Creation)
  - [Command](#Command)
    - [Get Version](#Get-Version)
    - [Command Help](#Command-Help)
    - [Command List](#Command-List)
    - [Command Creation Help](#Command-Creation-Help)
    - [Command Creation](#Command-Creation)
    - [Command Structure](#Command-Structure)
    - [Command Registration](#Command-Registration)
    - [Help Options](#Help-Options)
    - [Execute Command](#Execute-Command)
    - [Compile And Execute Commands](#Compile-And-Execute-Commands)
  - [Cache](#Cache)
    - [Global Cache](#Global-Cache)
    - [Redis Cache](#Redis-Cache)
    - [Memory Cache](#Memory-Cache)
    - [Disk Cache](#Disk-Cache)
  - [Event](#Event)
    - [Event Creation Help](#Event-Creation-Help)
    - [Event Creation](#Event-Creation)
  - [Listener](#Listener)
    - [Listener Creation Help](#Listener-Creation-Help)
    - [Listener Creation](#Listener-Creation)
  - [Queue](#Queue)
    - [Queue Creation Help](#Queue-Creation-Help)
    - [Queue Creation](#Queue-Creation)
    - [Queue Usage](#Queue-Usage)
  - [Publish Event](#Publish-Event)
    - [Event Test](#Event-Test)
  - [Event List](#Event-List)
    - [Event Listener List](#Event-Listener-List)
  - [Response](#Response)
    - [Response Success](#Response-Success)
      - [Response Success With Message](#Response-Success-With-Message)
      - [Response Success With Data](#Response-Success-With-Data)
    - [Response Error](#Response-Error)
      - [Response Error With Code](#Response-Error-With-Code)
      - [Response Error With Message](#Response-Error-With-Message)
      - [Response Error With Data](#Response-Error-With-Data)
  - [Log](#Log)
    - [Write Log](#Write-Log)
    - [Error Debug](#Error-Debug)
  - [Language Support](#Language-Support)
    - [Directory Configuration](#Directory-Configuration) 
    - [Ordinary Translation](#Ordinary-Translation) 
    - [Template Translation](#Template-Translation) 
    - [Add Language Support](#Add-Language-Support) 
  - [Provider Service](#Provider-Service)
    - [Provider Service Creation](#Provider-Service-Creation)
  - [Facade](#Facade)
    - [Facade Creation](#Facade-Creation)
    - [Facade Usage](#Facade-Usage)
  - [Database](#Database)
    - [Database Configuration](#Database-Configuration)
    - [Database Connection](#Database-Connection)
    - [Database Search](#Database-Search)
  - [Swagger Documents](#Swagger-Documents)

# Project Introduction
> A lightweight framework developed based on the Golang language framework `Go Gin`, out of the box, inspired by mainstream PHP frameworks such as `Laravel` and `ThinkPHP`. The project architecture directory has a clear hierarchy, which is a blessing for beginners. The framework integrates `jwt`, `log`、 `middleware`, `cache`, `validator`, `event`, `routing`, `queue(kafka、rabbitmq)`、 `redis`、 `Command` and other technologies. support multiple languages, simple to develop and easy to use, convenient for extension.
## Project Address
- Github: https://github.com/dsxwk/gin-admin.git
- Gitee: https://gitee.com/dsxwk/gin-admin.git

## Introduction to the Gin Framework
> Gin is a web framework written in Go language. It has the characteristics of simplicity, speed, and efficiency, and is widely used in Go language web development.

## Features of Gin Framework
- Fast: The Gin framework is based on the standard library net/http, using goroutines and channels to implement asynchronous processing and improve performance.
- Simple: The Gin framework provides a range of APIs and middleware, enabling developers to quickly build web applications.
- Efficient: The Gin framework uses sync. Pool to cache objects, reducing memory allocation and release, and improving performance.
> Golang Gin is a lightweight and efficient Golang web framework. It has the characteristics of high performance, ease of use, and flexibility, and is widely used in the development of various web applications.

# License
- 📘 Open source version: Following AGPL-3.0, for learning, research, and non-commercial use only.
- 💼 Commercial version: If closed source or commercial use is required, please contact the author 📧   [ 25076778@qq.com ]Obtain commercial authorization.

# Version History
> - Latest Version: v2.0.0
> - [Version update detailed record](VersionHistoryEn.md)

# Installation Instructions
> The project is developed based on Golang version 1.25.2, and there may be version differences in lower versions. It is recommended that the version be greater than or equal to 1.25.2.
## Clone Project
```bash
$ git clone https://github.com/dsxwk/gin.git
$ cd gin
```
## Initialize Go Environment And Dependencies
### Method One
```bash
$ go env -w GOPROXY=https://goproxy.cn,direct
$ go generate ./...
```
### Method Two
```bash
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
# $ go get -u
$ go mod tidy
# $ go mod download
$ go mod vendor
```
## Start
```bash
$ go run main.go
```
### Use Air Hot Update
```bash
$ go install github.com/air-verse/air@latest
$ air
```

## Compile
### Compile Project
```bash
$ go build main.go
$ ./main
```

### Compile Command
```bash
$ go build ./cmd/cli.go
$ ./cli demo:command --args=11

Excute Command: demo:command, Argument: 11
```

# Directory Structure
```
├── app                                 # Application
│   ├── command                         # Command
│   ├── controller                      # Controller
│   ├── event                           # Event
│   ├── facade                          # Facade
│   ├── listener                        # Listener
│   ├── middleware                      # Middleware
│   ├── model                           # Model
│   ├── provider                        # Provider
│   ├── queue                           # Queue
│   ├──├── kafka                        # Kafka
│   ├──├──├── consumer                  # Consumer
│   ├──├──├── producer                  # Producer
│   ├──├── rabbitmq                     # Rabbitmq
│   ├──├──├── consumer                  # Consumer
│   ├──├──├── producer                  # Producer
│   ├── request                         # Validator
│   ├── service                         # Service
├── bootstrap                           # Bootstrap 
├── cmd                                 # Command Script Tool
│   ├── cli.go                          # Entry File
├── common                              # Common Module
│   ├── base                            # Base
│   ├── ctxkey                          # Context Key
│   ├── errcode                         # Errcode
│   ├── flag                            # Flag
│   ├── response                        # Response
│   ├── template                        # Template
├── config                              # Config File
├── database                            # Database Test File 
├── docs                                # Swagger Doc
├── pkg                                 # Pakage
│   ├──├── cache                        # Cache
│   ├──├── cli                          # Command
│   ├──├── debugger                     # Debugger
│   ├──├── eventbus                     # Event Bus
│   ├──├── foundation                   # Providers
│   ├──├── http                         # Http Request
│   ├──├── lang                         # Language
│   ├──├── logger                       # Logger
│   ├──├── message                      # Message Event
│   ├──├── orm                          # Orm Tool
│   ├──├── queue                        # Queue
│   ├──├── time                         # Time Processing
├── public                              # Static Resources
├── router                              # Router
├── storage                             # Storage
│   ├── cache                           # Disk Cache
│   ├── logs                            # Logs
│   ├── locales                         # Translation
│   ├──├── en                           # English Translation
│   ├──├── zh                           # Chinese Translation
├── tests                               # Test Case
├── vendor                              # Vendor
├── .air.linux.toml                     # Air Configuration File
├── .air.toml                           # Air Configuration File
├── .gitignore                          # Gitignore
├── config.yaml                         # Default Configuration File
├── dev.config.yaml                     # Local Environment Configuration File
├── go.mod                              # go mod
├── LICENSE                             # LICENSE
├── main.go                             # Entry File
├── readme.md                           # English Document
├── readme_zh.md                        # Chinese Document
├── VersionHistoryEn.md                 # Version History English Document
└── VersionHistoryZn.md                 # Version History Chinese Document
```

# Instructions For Use
## Start Service
```bash
$ go run main.go
```
### Air Hot Update
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
应用:    gin
环境:    dev
端口:    8080
数据库:  gin
🌐 Address:    http://0.0.0.0:8080
👉 Swagger:    http://127.0.0.1:8080/swagger/index.html
👉 Test API:   http://127.0.0.1:8080/ping
 SUCCESS  Gin server started successfully!
```

## Configuration File
### Project Configuration
> `config.yaml` is the default configuration file and can be modified by oneself. `dev.config.yaml` corresponds to the local environment configuration, and environment variables can be configured through the following app.exe file to switch environments
> ```
> app:
>   env: dev # dev|testing|production dev=local-environment testing=test-environment production=production-environment
> ```

### Hot Update Configuration
> `.air.toml` is the default configuration file in Windows environment, and `.air.Linux.toml` is the default configuration file in Linux environment. You can modify it according to the overall needs of the project.

## Middleware
> `middleware`目录下为中间件目录, 可自行添加中间件, 并在`router/root.go`文件中注册中间件。
### Middleware Creation Help
```bash
$ go run ./cmd/cli.go make:middleware -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:middleware  Middleware Creation

Options:
  -f, --file  File Path, Expample: auth                        required:true
  -d, --desc  Description, Expample: Authorization-Middleware  required:false
```

### Middleware Creation
```bash
$ go run ./cmd/cli.go make:middleware --file=auth --desc=Authorization-Middleware
```

### Rate Limit Middleware
> The `middleware/rate_imit.go` file defines a global flow limiting middleware that supports global user interface flow limiting, IP interface flow limiting, and global flow limiting.
```go
package router

import (
    "gin/app/middleware"
    "github.com/gin-gonic/gin"
)

var rateLimitMiddleware middleware.RateLimit

// LoadRouters Load Routers
func LoadRouters(router *gin.Engine) {
    // Global Rate Limit
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

	// Specify interface current limit
    // User Rate Limit
    // r How many tokens are generated per second
    // burst Bucket capacity
    userGroup := router.Group("", rateLimitMiddleware.UserRateLimit(1, 1))
	r1 := userGroup.Group("")
	{
		r1.GET("/test1", func(c *gin.Context) {
			err := errcode.NewError(0, "user test1")
			response.Success(c, &err)
		})

		r1.GET("/test2", func(c *gin.Context) {
			err := errcode.NewError(0, "user test2")
			response.Success(c, &err)
		})
    }

    // Specify interface current limit
    // Ip Rate Limit
    // r How many tokens are generated per second
    // burst Bucket capacity
    ipGroup := router.Group("", rateLimitMiddleware.IpRateLimit(1, 1))
	r2 := ipGroup.Group("")
	{
		r2.GET("/test1", func(c *gin.Context) {
			err := errcode.NewError(0, "ip test1")
			response.Success(c, &err)
		})

		r2.GET("/test2", func(c *gin.Context) {
			err := errcode.NewError(0, "ip test2")
			response.Success(c, &err)
		})
    }
}
```

## Route
> The `router/root.go` file defines global routing rules that can be modified by oneself, and in general, they only need to be defaulted.
### Route Creation Help
```bash
$ go run ./cmd/cli.go make:router -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:router  Route Creation

Options:
  -f, --file  File Path, Expample: user                   required:true
  -d, --desc  Route Description, Expample: User-Routing   required:false
```

### Route Creation
```bash
$ go run ./cmd/cli.go make:router --file=user --desc=User-Routing
```
```go
package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// UserRouter User-Routing
type UserRouter struct{}

func init() {
	Register(&UserRouter{})
}

// RegisterRoutes Register-Route
func (r *UserRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		user v1.UserController
	)

	router := routerGroup.Group("api/v1")
	{
		// List
		router.GET("/user", user.List)
		// Create
		router.POST("/user", user.Create)
		// Update
		router.PUT("/user/:id", user.Update)
		// Delete
		router.DELETE("/user/:id", user.Delete)
		// Detail
		router.GET("/user/:id", user.Detail)
	}
}

// IsAuth Do you need authentication
func (r *UserRouter) IsAuth() bool {
	return true
}

```

### Route List
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
A total of 10 routes
```

## Controller
### Controller Creation Help
```bash
$ go run ./cmd/cli.go make:controller -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:controller  Controller Creation

Options:
  -f, --file      File Path, Example: v1/user       required:true
  -F, --function  Function Name, Example: list      required:false
  -m, --method    Request Method, Example: get      required:false
  -r, --router    Route Adress, Example: /user      required:false
  -d, --desc      Description, Example: Test-List   required:false
```

### Controller Creation
```bash
$ go run ./cmd/cli.go make:controller --file=v1/test --router=/test --method=get --desc=Test-List --function=list
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

// List Test-List
// @Router /test [get]
func (s *TestController) List(c *gin.Context) {
    // Define your function here
    s.Success(c, errcode.Success().WithMsg("Test Msg").WithData([]string{}))
}
```
## Model
### Model Creation Help
```bash
$ go run ./cmd/cli.go make:model -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:model  Model Creation

Options:
  -t, --table       Table Name, Example: user or user,menu     required:true
  -p, --path        Output Directory, Example: api/user        required:false
  -c, --camel       Is it a camel hump field, Example: true    required:false
  -C, --connection  Database Connection                        required:false
```

### Model Creation
> Support the creation of multiple model files simultaneously. If multiple model files need to be created, please separate the table name parameters of the descendants with commas, such as: user, menu
```bash
$ go run ./cmd/cli.go make:model --table=user,menu --path=api/user --camel=true
```
```go
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package user

import "gin/app/model"

const TableNameUser = "user"

// User User-Table
type User struct {
	ID        int64            `gorm:"column:id;type:int(10) unsigned;primaryKey;autoIncrement:true;comment:ID" json:"id"`                       // ID
	Avatar    string           `gorm:"column:avatar;type:varchar(255);not null;comment:avatar" json:"avatar"`                                    // avatar
	Username  string           `gorm:"column:username;type:varchar(10);not null;comment:username" json:"username"`                               // username
	FullName  string           `gorm:"column:full_name;type:varchar(20);not null;comment:fullname" json:"fullName"`                              // fullName
	Email     string           `gorm:"column:email;type:varchar(50);not null;comment:email" json:"email"`                                        // email
	Password  string           `gorm:"column:password;type:varchar(255);not null;comment:password" json:"password"`                              // password
	Nickname  string           `gorm:"column:nickname;type:varchar(50);not null;comment:nickname" json:"nickname"`                               // nickname
	Gender    int64            `gorm:"column:gender;type:tinyint(1) unsigned;not null;comment:gender 1=male 2=female" json:"gender"`             // gender 1=male 2=female
	Age       int64            `gorm:"column:age;type:int(10) unsigned;not null;comment:age" json:"age"`                                         // age
	Status    int64            `gorm:"column:status;type:tinyint(3) unsigned;not null;default:1;comment:state 1=enable 2=disable" json:"status"` // state 1=enable 2=disable
	CreatedAt *model.DateTime  `gorm:"column:created_at;type:datetime;comment:Creation Time" json:"createdAt"`                                   // Creation Time
	UpdatedAt *model.DateTime  `gorm:"column:updated_at;type:datetime;comment:Update Time" json:"updatedAt"`                                     // Update Time
	DeletedAt *model.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:Delete Time" json:"deletedAt" swaggerignore:"true"`                // Delete Time
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
```

## ORM Dynamic Filtering
> By passing the `query` | `body` parameter `__search` through `post` or `get`, dynamically specify the query criteria based on the list fields. The `__search` type is `map[string]interface{}`, for example:__ search={"and":[{"username":"test"},{"age":18}]}, __search={"or":[{"username":"test"},{"age":18}]}.  support or、and、in、not in、between、not between、like、left like、right like、is not null、is null、gt、gte、lt、lte、exist、not exist、json_contains、json_extract Wait for conditions, case insensitive The parameter supports two modes: `{'username': 'admin'}` or `{'username': ['like', 'admin']}`. When the field name is a keyword of the 'mysql where' condition, SQL statements will be automatically constructed based on the condition
### OR Condition Query
```http
GET /api/v1/user?__search={"or":[{"username":"test"},{"age":18}]} // {"or":[{"username":["=", "test"]},{"age":["=", 18]}]}
```
```sql
SELECT * FROM `user` WHERE (username = 'test' OR age = 18)
```

### AND Condition Query
```http
GET /api/v1/user?__search={"and":[{"username":"test"},{"age":18}]} // {"and":[{"username":["=", "test"]},{"age":["=", 18]}]}
```
```sql
SELECT * FROM `user` WHERE (username = 'test' AND age = 18)
```

### JSON Field Query
```http
GET /api/v1/menu?__search={"or":[{"and":[{"createdAt":[">","2025-01-01"]},{"createdAt":["<","2026-01-01"]},{"name":""},{"$.meta.icon":["=","ele-Collection"]}]}]}
```
```sql
 SELECT * FROM `menu` WHERE ((((menu.created_at > '2025-01-01') AND (menu.created_at < '2026-01-01') AND (menu.name = '') AND (JSON_EXTRACT(meta, '$.icon') = 'ele-Collection'))))
```

### Complex Condition Query
```http
GET /api/v1/user?__search={"or":[{"and":[{"createdAt":[">","2025-01-01"]},{"createdAt":["<","2026-01-01"]},{"not exist":{"userRoles.name":"admin"}}]},{"username":"admin"}]}
```
```sql
 SELECT * FROM `user` WHERE ((((user.created_at > '2025-01-01') AND (user.created_at < '2026-01-01') AND (NOT EXISTS (SELECT 1 FROM user_roles WHERE user_roles.user_id = user.id AND user_roles.name = 'admin'))) OR (user.username = 'admin')))
```

## Form Validation
### Validator Creation Help
```bash
$ go run ./cmd/cli.go make:request -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:request  Validator Creation

Options:
  -f, --file  File Path, Example: user                         required:true
  -d, --desc  Description, Example: user-request-validation    required:false
```

### Validator Creation
```bash
$ go run ./cmd/cli.go make:request --file=user --desc=user-request-validation
```
```go
package request

import (
    "errors"
    "github.com/gookit/validate"
)

// User User-Request-Validation
type User struct {
    PageListValidate
}

// Validate Request-Validation
func (s User) Validate(data User, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}

// ConfigValidation Configuration-Validation
// - Define validation scenes
// - You can also add verification settings
func (s User) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"create": []string{},
		"update": []string{"ID"},
		"detail": []string{"ID"},
		"delete": []string{"ID"},
	})
}

// Messages Validator-Error-Message
func (s User) Messages() map[string]string {
	return validate.MS{
		"required":    "Field {field} Required",
		"int":         "Field {field} Must be an integer",
		"Page.gt":     "Field {field} Must be greater than 0",
		"PageSize.gt": "Field {field} Must be greater than 0",
	}
}

// Translates Field-Translation
func (s User) Translates() map[string]string {
	return validate.MS{
		"Page":     "Page",
		"PageSize": "Page Size",
		"ID":       "ID",
	}
}
```

### Validator Rules
> For more rules, please refer to [gookit/validate](https://github.com/gookit/validate)
```go
package request

// UserCreate User-Create-Validation
type UserCreate struct {
	Username string `json:"username" validate:"required" label:"username"`
	FullName string `json:"fullName" validate:"required" label:"fullname"`
	Nickname string `json:"nickname" validate:"required" label:"nickname"`
	Gender   int    `json:"gender" validate:"required|int" label:"gender"`
	Password string `json:"password" validate:"required" label:"password"`
}

// UserUpdate User-Update-Validation
type UserUpdate struct {
    ID int64 `json:"id" validate:"required|int|gt:0" label:"ID"`
    Username string `json:"username" validate:"required" label:"username"`
    FullName string `json:"fullName" validate:"required" label:"fullname"`
    Nickname string `json:"nickname" validate:"required" label:"nickname"`
    Gender   int    `json:"gender" validate:"required|int" label:"gender"`
    Password string `json:"password" validate:"required" label:"password"`
}

// UserDetail User-Detail-Validation
type UserDetail struct {
    ID int64 `json:"id" validate:"required|int|gt:0" label:"ID"`
}

// User User-Request-Validation
type User struct {
    ID int64 `json:"id" validate:"required|int|gt:0" label:"ID"`
    Username string `json:"username" validate:"required" label:"username"`
    FullName string `json:"fullName" validate:"required" label:"fullname"`
    Nickname string `json:"nickname" validate:"required" label:"nickname"`
    Gender   int    `json:"gender" validate:"required|int" label:"gender"`
    Password string `json:"password" validate:"required" label:"password"`
	PageListValidate
}
```

### Validator Scenes
```go
package request

// ConfigValidation Configuration-Validation
// - Define validation scenes
// - You can also add verification settings
func (s User) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		// List
		"List": []string{
			"PageListValidate.Page",
			"PageListValidate.PageSize",
		},
		// Create
		"Create": []string{
			"Username",
			"FullName",
			"Nickname",
			"Gender",
			"Password",
		},
		// Update
		"Update": []string{
			"ID",
			"Username",
			"FullName",
			"Nickname",
			"Gender",
		},
		// Detail
		"Detail": []string{
			"ID",
		},
		// Delete
		"Delete": []string{
			"ID",
		},
	})
}
```

### Prompt Message
```go
package request

// Messages Validator-Error-Message
func (s User) Messages() map[string]string {
    return validate.MS{
        "required":                     "Field {field} Required",
        "int":                          "Field {field} Must be an integer",
        "PageListValidate.Page.gt":     "Field {field} Must be greater than 0",
        "PageListValidate.PageSize.gt": "Field {field} Must be greater than 0",
    }
}
```

### Field Translation
```go
package request

// Translates Field-Translation
func (s User) Translates() map[string]string {
	return validate.MS{
		"Page":     "Page",
		"PageSize": "Page Size",
		"ID":       "ID",
		"Username": "Username",
		"FullName": "Fullname",
		"Nickname": "Nickname",
		"Gender":   "Gender",
		"Password": "Password",
	}
}
```

### Custom Validation
#### Global Rules
> Global rules only need to be defined in the entry file `main.go`,applicable to all validators, without the need for repeated definitions.
```go
package main

import (
  "github.com/gookit/validate"
)

// Register during initialization
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

#### Local Rules
```go
package request

// ValidateIsEven Define local rule methods (naming convention: Validate<rule name>)
func (s User) ValidateIsEven(val any) bool {
    num := val.(int)
    return num%2 == 0
}
```

#### Temporary Rules
```go
package request

// Validate Request-Validation
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

#### Validator Usage
```go
package request

type User struct {
    Age int `json:"gender" validate:"required|is_even" label:"age"`
}
```

#### Used In The Controller
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

// List User-List
// @Tags User
// @Summary List
// @Description User-List
// @Param token header string true "Authentication Token"
// @Param page query string true "Page"
// @Param pageSize query string true "Page Size"
// @Success 200 {object} errcode.SuccessResponse{data=request.PageData{list=[]model.User}} "Login Successful"
// @Failure 400 {object} errcode.ArgsErrorResponse "Argument Error"
// @Failure 500 {object} errcode.SystemErrorResponse "System Error"
// @Router /api/v1/user [get]
func (s *UserController) List(c *gin.Context) {
  var (
    ctx    = c.Request.Context()
    req    request.User
    search request.Search
  )

  s.service.WithContext(ctx)

  err := c.ShouldBind(&search)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(err.Error()))
    return
  }

  err = c.ShouldBind(&req)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(err.Error()))
    return
  }

  // Validator
  err = req.Validate(req, "List")
  if err != nil {
    s.Error(c, errcode.ArgsError().WithMsg(err.Error()))
    return
  }

  res, err := s.service.List(req, search.Search)
  if err != nil {
    s.Error(c, errcode.SystemError().WithMsg(facade.Lang.T(ctx, err.Error(), nil)))
    return
  }

  s.Success(c, errcode.Success().WithData(res))
}
```

## Service
### Service Creation Help
```bash
$ go run ./cmd/cli.go make:service -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:service  Service Creation

Options:
  -f, --file      File Path, Example: v1/user      required:true
  -F, --function  Function Name, Example: list     required:false
  -d, --desc      Description, Example: list       required:false
exit status 3
```

### Service Creation
```bash
$ go run ./cmd/cli.go make:service -f=user --function=list --desc="list"
```

## Command
### Get Version
```bash
$ go run ./cmd/cli.go --version # -v
Gin Cli v2.0.0
```

### Command Help
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
  make:command     Command Creation
  make:controller  Controller Creation
  make:event       Event Creation
  make:facade      Facade Creation
  make:listener    Listener Creation
  make:middleware  Middleware Creation
  make:migration   Generate database migration template
  make:model       Model Creation
  make:model-old   Model Creation old
  make:provider    Producer Creation
  make:queue       Queue Creation(Kafka/RabbitMQ)
  make:request     Request Creation
  make:router      Router Creation
  make:seed        Generate database seeder template
  make:service     Service Creation
producer:
  producer:list    Producer List
route:
  route:list       Route List

Options:
  -f, --format     The output format (txt, json) [default: txt]
  -h, --help       Display help for the given command
  -v, --version    Display CLI version
```

### Command List
```bash
$ go run ./cmd/cli.go --format=json # -f=json
{
  "commands": [
    {
      "description": "Consumer List",
      "name": "consumer:list"
    },
    {
      "description": "Database migrate",
      "name": "db:migrate"
    },
    {
      "description": "Database rollback",
      "name": "db:rollback"
    },
    {
      "description": "Database seed",
      "name": "db:seed"
    },
    {
      "description": "Demo test",
      "name": "demo:command"
    },
    {
      "description": "Event List",
      "name": "event:list"
    },
    {
      "description": "Listener List",
      "name": "listener:list"
    },
    {
      "description": "Command Creation",
      "name": "make:command"
    },
    {
      "description": "Controller Creation",
      "name": "make:controller"
    },
    {
      "description": "Event Creation",
      "name": "make:event"
    },
    {
      "description": "Facade Creation",
      "name": "make:facade"
    },
    {
      "description": "Listener Creation",
      "name": "make:listener"
    },
    {
      "description": "Middleware Creation",
      "name": "make:middleware"
    },
    {
      "description": "Generate database migration template",
      "name": "make:migration"
    },
    {
      "description": "Model Creation",
      "name": "make:model"
    },
    {
      "description": "Model Creation old",
      "name": "make:model-old"
    },
    {
      "description": "Provider Creation",
      "name": "make:provider"
    },
    {
      "description": "Queue Creation(Kafka/RabbitMQ)",
      "name": "make:queue"
    },
    {
      "description": "Request Creation",
      "name": "make:request"
    },
    {
      "description": "Route Creation",
      "name": "make:router"
    },
    {
      "description": "Generate database seeder template",
      "name": "make:seed"
    },
    {
      "description": "Service Creation",
      "name": "make:service"
    },
    {
      "description": "Producer List",
      "name": "producer:list"
    },
    {
      "description": "Route List",
      "name": "route:list"
    }
  ],
  "version": "Gin Cli v2.0.0"
}
```

## Command Creation Help
```bash
$ go run ./cmd/cli.go make:command -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
make:command  Command Creation

Options:
  -f, --file  File Path, Example: cronjob/demo     required:true
  -n, --name  Command Name, Example: demo-test     required:false
  -d, --desc  Description, Example: command-desc   required:false
```

## Command Creation
```bash
$ go run ./cmd/cli.go make:command --file=cronjob/demo --name=demo-test --desc=command-desc
```

## Command Structure
> After generating the command, appropriate values should be defined for the ` Name() ` and ` Descript() ` functions. These properties will be used when displaying the command list. The `Name()` function also allows you to define the expected input value for the command. It will call the `Execute()` function when executing the command. You can put the command logic in this method. Let's take a look at an example command.
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
            "Example Argument, Example: arg1",
            true,
        },
    }
}

func (m *DemoCommand) Execute(args []string) {
    values := m.ParseFlags(m.Name(), args, m.Help())
    color.Green("Execute Command: %s %s", m.Name(), m.FormatArgs(values))
}

func init() {
	cli.Register(&DemoCommand{})
}

```

## Command Registration
> `cli. go` registers all commands in the `command` package under the `gin/app/command` directory by default. If the command you registered is not `command` package, you can add the path to import the package in `./cmd/imports/import.go`.
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

## Help Options
> Command option parameters are defined using the `base. CommandOption` structure. The `base. CommandOption` struct contains two attributes: `Flag` and `Description`. The `Flag` attribute is used to define the flag of command options, which can be a short flag (such as `- a `) or a long flag (such as `--args`). The `Description` attribute is used to define the description of command options. The `base. CommandOption` struct also contains a `Required` attribute that specifies whether a command option is required. At the same time, this method supports the console `--help` parameter and automatically generates help information.
```go
func (m *DemoCommand) Help() []base.CommandOption {
	return []base.CommandOption{
        {
            base.Flag{
                Short: "a",
                Long:  "args",
            },
            "Example Argument, Example: arg1",
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
  -a, --args  Example Argument, Example: arg1  required:true
```

## Execute Command
```bash
$ go run ./cmd/cli.go demo:command --args=arg1
 SUCCESS  Excute Command: demo:command --args=arg1
```

## Compile And Execute Commands
```bash
$ go build ./cmd/cli.go
$ ./cli demo:command --args=arg1
```

# Cache
> With `memory` as the default cache driver and support for custom extensions. By default, it supports three modes: `Memory cache`, `Redis cache`, and `Disk cache`. It can use global cache or any cache separately. The global cache only integrates the common methods of `Set`, `Get`, `Delete`, and `Expire` by default. If you need to use more, you can use them separately, or you can integrate them yourself.
## Global Cache
> The configuration of global cache can be switched through the `cache.driver` configuration in the `yaml` configuration file, or dynamically switched.
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

func (s *TestController) Test()  {
    // Set Set-Cache	
    key := "test_key"
    value := "test_value"
    cache := facade.Cache.Store()
    cache = facade.Cache.Store("redis")
    err := cache.Set(key, value, time.Second*10)
	if err != nil {
	    // Handle error	
    }
	
    // Get Get-Cache
    key := "test_key"
    value := "test_value"
    result, ok := cache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete Delete-Cache
	key := "test_key"
	err := cache.Delete(key)
	if err != nil {
        // Handle error	
    }
	
	// Expire Get-Cache-Expire
	key := "test_key"
    val, expireAt, ok, err := cache.Expire(key)
	if err != nil {
	    // Handle error
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
}
```

## Redis Cache
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

func (s *TestController) Test()  {
    // Set Set-Cache	
    key := "test_key"
    value := "test_value"
    redisCache := facade.Cache.Store("redis")
    err := redisCache.Set(key, value, time.Second*10)
	if err != nil {
	    // Handle error	
    }
	
    // Get Get-Cache
    key := "test_key"
    value := "test_value"
    result, ok := redisCache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete Delete-Cache
	key := "test_key"
	err := redisCache.Delete(key)
	if err != nil {
        // Handle error	
    }
	
	// Expire Get-Cache-Expire
	key := "test_key"
    val, expireAt, ok, err := redisCache.Expire(key)
	if err != nil {
	    // Handle error
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
	
	// ... Other
}
```

## Memory Cache
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

func (s *TestController) Test()  {
    // Set Set-Cache	
    key := "test_key"
    value := "test_value"
    memoryCache := facade.Cache.Store("memory")
    err := memoryCache.Set(key, value, time.Second*10)
	if err != nil {
	    // Handle error	
    }
	
    // Get Get-Cache
    key := "test_key"
    value := "test_value"
    result, ok := memoryCache.Get(key)
	if ok {
	    println(result) // test_value	
    }
	
	// Delete Delete-Cache
	key := "test_key"
	err := memoryCache.Delete(key)
	if err != nil {
        // Handle error	
    }
	
	// Expire Get-Cache-Expire
	key := "test_key"
    val, expireAt, ok, err := memoryCache.Expire(key)
	if err != nil {
	    // Handle error
    }
	if ok {
      fmt.Println(val) // test_value
      fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
	
	// ... Other
}
```

## Disk Cache
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
    // Set Set-Cache	
    key := "test_key"
    value := "test_value"
    diskCache := facade.Cache.Store("disk")
    err := diskCache.Set(key, value, time.Second*10)
    if err != nil {
        // Handle error	
    }
    
    // Get Get-Cache
    key := "test_key"
    value := "test_value"
    result, ok := diskCache.Get(key)
    if ok {
        println(result) // test_value	
    }
    
    // Delete Delete-Cache
    key := "test_key"
    err := diskCache.Delete(key)
    if err != nil {
        // Handle error	
    }
    
    // Expire Get-Cache-Expire
    key := "test_key"
    val, expireAt, ok, err := diskCache.Expire(key)
    if err != nil {
        // Handle error
    }
    if ok {
        fmt.Println(val) // test_value
        fmt.Printf("ExpireAt: %v\n", expireAt) // ExpireAt: 2025-10-28 11:23:38.7416956 +0800 CST
    }
    
    // ... Other
}	
```

# Event
## Event Creation Help
```bash
$ go run ./cmd/cli.go make:event -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:event  Event Creation

Options:
  -f, --file  File Path, Example: login/test             required:true
  -n, --name  Event Name, Example: test-event            required:false
  -d, --desc  Event Description, Example: test-event     required:false
```

## Event Creation
```bash
$ go run ./cmd/cli.go make:event -f=user_login -n='user.login' -d=user-login-event
```
```go
package event

// UserLoginEvent event-data
type UserLoginEvent struct {
	UserId   int64
	Username string
}

// Name event-name
func (u UserLoginEvent) Name() string {
	return "user.login"
}

// Description event-description
func (u UserLoginEvent) Description() string {
	return "User Login Event"
}

```

# Listener
## Listener Creation Help
```bash
$ go run ./cmd/cli.go make:listener -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:listener  Listener Creation

Options:
  -f, --file   File Path, Example: login/test   required:true
  -e, --event  Event Data, Example: UserLogin   required:true
```

## Listener Creation
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
	fmt.Printf("Recieved Event: %s Event Description: %s Event Data: %s, Time: %s\n", ev.Name(), ev.Description(), data, time.Now().Format("2006-01-02 15:04:05"))
}

func init() {
	eventbus.Register(&UserLoginListener{}, event.UserLoginEvent{})
}

```

# Queue
> Executing the queue creation command will create both consumers and producers based on the queue type. For example, Kafka will create Kafka consumers and producers, while RabbitMQ will create RabbitMQ consumers and producers You only need to improve the `Handle` method among consumers to enhance your business logic, supporting automatic error retries and delayed queues
## Queue Creation Help
```bash
$ go run ./cmd/cli.go make:queue -h # --help
Gin Cli v2.0.0

Usage:
  cli [command] [options]

Command:
  make:queue  Queue Creation

Options:
  -t, --type      Queue Type, Example: kafka or rabbitmq     required:true
  -n, --name      Queue File Name, Example: order_create     required:true
  -d, --isDelay   Is Delay Queue, Example: true or false     required:false
  -T, --topic     Queue Topic, Example: kafka_demo           required:false
  -k, --key       Message Key, Example: kafka_demo           required:false
  -g, --group     Goup, Example: kafka_demo                  required:false
  -q, --queue     Queue Name, Example: rabbitmq_demo         required:false
  -e, --exchange  Exchange, Example: rabbitmq_demo           required:false
  -r, --routing   Routing Key, Example: rabbitmq_demo        required:false
  -R, --retry     Retry Times, Example: 3                    required:false
  -m, --delayMs   Delay In Milliseconds, Example: 10000      required:false
```

## Queue Creation
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

// RabbitmqDemoConsumer RabbitMQ Consumer
type RabbitmqDemoConsumer struct {
  *base.RabbitmqConsumer
}

// NewRabbitmqDemoConsumer Consumer Creation
func NewRabbitmqDemoConsumer() *RabbitmqDemoConsumer {
  cfg := facade.Config.Get()
  log := facade.Log.Logger()
  bus := facade.Message.GetBus()

  // RabbitMQ Connection Creation
  mq, err := base.NewRabbitMQ(cfg, log, bus)
  if err != nil {
    log.Error(pkg.Sprintf("RabbitMQ Connection Failed: %v", err))
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

func (c *RabbitmqDemoConsumer) Name() string {
  return "rabbitmq_demo"
}

func (c *RabbitmqDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
  c.RabbitmqConsumer.Start(c)
  log.Info(pkg.Sprintf("RabbitMQ Consumer Start Successed: %s", c.Name()))
  return nil
}

func (c *RabbitmqDemoConsumer) Stop() error {
  return c.RabbitmqConsumer.Stop()
}

func (c *RabbitmqDemoConsumer) Enabled(cfg *config.Config) bool {
  return cfg.Rabbitmq.Enabled
}

func (c *RabbitmqDemoConsumer) Status() queue.ConsumerStatus {
  return c.RabbitmqConsumer.Status()
}

// Handle Process business logic
func (c *RabbitmqDemoConsumer) Handle(msg string) error {
  facade.Log.Info(pkg.Sprintf("RabbitMq Received Msg: %s", msg))
  // todo Process business logic
  return nil
}

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

// RabbitmqDemoProducer RabbitMQ Producer
type RabbitmqDemoProducer struct {
  *base.RabbitmqProducer
}

// NewRabbitmqDemoProducer Producer Creation
func NewRabbitmqDemoProducer() *RabbitmqDemoProducer {
  cfg := facade.Config.Get()
  log := facade.Log.Logger()
  bus := facade.Message.GetBus()

  mq, err := base.NewRabbitMQ(cfg, log, bus)
  if err != nil {
    log.Error(pkg.Sprintf("RabbitMQ Connection Failed: %v", err))
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

## Queue Usage
> When consumers start a project, it is automatically registered in the container for unlimited additional launches, and producers can use it directly by initializing the storefront.
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
    // Get Producer
    producer := facade.Queue.Producer("rabbitmq_demo")
	_ = producer.Publish(ctx, []byte(`{"orderId":111, "message":"message 111"}`))
}
```

# Publish Event
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
)

type LoginController struct {
    base.BaseController
    service service.LoginService
}

// Token token-info
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

// Login login
// @Tags login
// @Summary login
// @Description User login
// @Accept json
// @Produce json
// @Param data body request.UserLogin true "Login Argument"
// @Success 200 {object} errcode.SuccessResponse{data=LoginResponse} "Success"
// @Failure 400 {object} errcode.ArgsErrorResponse "Argument Error"
// @Failure 500 {object} errcode.SystemErrorResponse "System Error"
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

  // Validate
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

  // Publish Event
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

## Event Test
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

收到事件: user.login Event Description: User Login Event Event Data: {"UserId":1,"Username":"admin"}, Time: 2025-11-04 15:32:12
```

# Event List
```bash
$ go run ./cmd/cli.go event:list

user.login 用户登录事件
```

## Event Listener List
```bash
$ go run ./cmd/cli.go listener:list

==== Currently registered events ====
Event: user.login
Description: User Login Event
Listener:
  - *listener.TestListener
  - *listener.UserLoginListener
----------------------
```

# Response
## Response Success
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

###  Response Success With Message
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

### Response Success With Data
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

## Response Error
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

### Response Error With Code
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

### Response Error With Message
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

### Response Error With Data
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

# Log
> Use the `zap` package to implement logging. The storage path for log files is `storage/logs`, and the default log level is `debug`. When the error code returned is not 0, it automatically records log TraceId, stack, SQL, HTTP, Redis, and other call information. Logging can also be directly called to automatically record debugging information. Does `log.access` in the configuration file `yaml` support automatic recording of request logs? If enabled, it will automatically record request logs.
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

## Write Log
> Encapsulated in the facade, the log level supports debug, info, warn, error, dpanic, panic, and fatal, with the default being `debug`.
```go
package v1

import (
    "gin/app/facade"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    facade.Log.Error("System Error")
}
```

## Error Debug
> When using public return errors and calling the Withdebugger() method, it will automatically record log TraceId, stack, SQL, HTTP, Redis, and other call information. Debugging can be done based on debug and trace stack information. The log file storage path is' storage/logs'.
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

# Language Support
> Multilingualism has been integrated into the facade and provider, supporting both `zh` and `en` languages, and supporting custom extensions. Language transmission defaults to transmitting the `Accept-Language` parameter in the `header`, such as `zh` or `en`, which is not case sensitive and does not pass the default language as `zh`.
## Directory Configuration
> The storage path for translation files is `storage/scales`, the default language is `zh`, and multiple languages are separated by commas. Languages are stored in the corresponding language directory without distinguishing between subdirectories. For example, Chinese is stored in `storage/scales/zh` and can support `json` and `yaml` format files in any directory.
```yaml
# Translation Configuration
i18n:
  dir: "storage/locales" # Translation file storage path
  lang: "zh,en" # Default language, multiple languages separated by commas
```

## Ordinary Translation
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
	fmt.Println(trans) // Output: 用户名, English Output: Username
}
```

## Template Translation
> Template translation is supported in the translation file, such as `{{. name}}`, using `map[string]interface{}` to pass parameters.
```json
[
  {
    "id": "login.success",
    "translation": "{{.name}},Login Success"
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
	fmt.Println(trans) // Output: admin,登录成功 English Output: admin,Login Success
}
```

## Add Language Support
> Add the corresponding language directory, such as `en`, in the `storage/scales` directory, and then add a translation file in the directory. The translation file supports `json` and `yaml` formats, with `id` as the unique identifier and `translation` as the translation content. Any number of translation contents can be added to the translation file. The configuration language support requires adjusting the `i18n.lang` parameter in the configuration file.
```yaml
# Translation Configuration
i18n:
  dir: "storage/locales" # Translation file storage path
  lang: "zh,en" # Default language, multiple languages separated by commas
```

# Provider Service
> The provider service will automatically load the registration upon startup and release it upon shutdown.
## Provider Service Creation
> 同模型、控制器等使用命令行创建,具体参考之前文档。

# Facade
## Facade Creation
> Create models, controllers, etc. using the command line, refer to the previous documentation for details.

## Facade Usage
> The project integrates features such as logs, databases, caches, and throttling by Facade. Currently, cache is used as an example. The binding of context to databases, caches, HTTP requests, and queues will be recorded in the debugging log.
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
    redisCache := facade.Cache.Store('redis')   // OR facade.Cache.Redis()
    // Bind context to cache
    redisCache = redisCache.WithContext(ctx)
    memoryCache := facade.Cache.Store('memory') // OR facade.Cache.Memory()
    diskCache := facade.Cache.Store('disk')     // OR facade.Cache.Disk()
    // Other facade usage ...
}
```

# Database
> The database is initialized through a container and bound to the context through middleware, so that database instances can be obtained wherever there is context. You can also obtain database instances separately. By default, MySQL, pgSQL, SQLite, and SQLSRV are integrated, and the default database can be configured and the database connection can be specified through the Connection method.
## Database Configuration
```yaml
# Database
databases:
  db-connection: mysql # Default database connection
  # Slow query time (ms) exceeding this time will be recorded in the log
  slow-query-duration: 3000ms # 3 Second(time.Duration)

# Mysql Database
mysql:
  driver: mysql
  # host: "username:password@tcp(127.0.0.1:3306)/databaseName?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: gin
  # Slow query time (ms) exceeding this time will be recorded in the log
  slow-query-duration: 3000ms # 3 Second(time.Duration)

# Postgresql Database
pgsql:
  driver: pgsql
  host: 127.0.0.1
  port: 5432
  username: testuser
  password: 123456
  database: testdb
  # Slow query time (ms) exceeding this time will be recorded in the log
  slow-query-duration: 3000ms # 3 Second(time.Duration)

# sqlite Database
sqlite:
  driver: sqlite
  path: storage/data/gin.db
  # Slow query time (ms) exceeding this time will be recorded in the log
  slow-query-duration: 3000ms # 3 Second(time.Duration)

# sqlsrv Database
sqlsrv:
  driver: sqlsrv
  host: 127.0.0.1
  port: 1433
  username: root
  password: root
  database: gin
  # Slow query time (ms) exceeding this time will be recorded in the log
  slow-query-duration: 3000ms # 3 Second(time.Duration)
```

## Database Connection
> The use of context is not mandatory. If the context is not bound, SQL records will not be recorded in the log.
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
  // Default Connection
  db := facade.DB.Connection()
  // Using context
  db1 := facade.DB.Connection().WithContext(ctx)
  // Connection pgsql
  db2 := facade.DB.Connection("pgsql").WithContext(ctx)
  // Connection sqlsrv
  db3 := facade.DB.Connection("sqlsrv").WithContext(ctx)
  // todo ...
}
```

## Database Search
> Use in conjunction with the ORM dynamic filtering example in the document.
```go
package controller

import (
    "gin/app/facade"
    "gin/app/model"
    "gin/app/request"
    "github.com/gin-gonic/gin"
)

type TestController struct {
    base.BaseController
}

func (s *TestController) Test(c *gin.Context) {
    var (
        ctx    = c.Request.Context()
		search request.Search
        m      []model.User
		db     = facade.DB.Connection().WithContext(ctx)
	)

    err := c.ShouldBind(&search)
    if err != nil {
        s.Error(c, errcode.SystemError().WithMsg(err.Error()))
        return
    }

    db = db.Model(&model.User{})
	
    if search != nil {
        whereSql, args, err := orm.BuildCondition(search, db, model.User{})
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

# Swagger Documents
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

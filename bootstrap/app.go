package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"gin/app/facade"
	_ "gin/app/listener"
	_ "gin/app/provider"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/foundation"
	"gin/router"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-runewidth"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

// App 应用结构
type App struct {
	Engine *gin.Engine
}

// Init 初始化应用
func Init() error {
	// 初始化门面系统
	facade.Init()
	// 创建应用实例
	app := foundation.GetApp()
	// 注册应用到门面
	facade.Register("app", app)

	// 启动应用(加载所有providers)
	return app.Boot()
}

// NewApp 创建应用实例
func NewApp() *App {
	if err := Init(); err != nil {
		flag.Errorf("初始化应用失败: %v", err)
		os.Exit(1)
	}

	return &App{
		Engine: setupEngine(),
	}
}

// setupEngine 配置Gin引擎
func setupEngine() *gin.Engine {
	r := gin.New()

	conf := facade.Config.Get()

	// 设置运行模式
	gin.SetMode(conf.App.Mode)

	// 非生产环境允许所有代理
	if conf.App.Env != "production" {
		_ = r.SetTrustedProxies(nil)
	}

	// 设置文件上传最大内存
	r.MaxMultipartMemory = 90 << 20

	// 加载路由
	router.LoadRouters(r)

	return r
}

func (a *App) Run() {
	// 系统优化设置50：更积极回收
	debug.SetGCPercent(50)

	conf := facade.Config.Get()
	if conf == nil {
		flag.Errorf("配置未加载,无法启动服务")
		os.Exit(1)
	}

	// 启动提示
	a.printStartupInfo(conf)

	// 启动http服务
	srv := a.startHttpServer(conf)

	// 优雅关闭
	a.gracefulShutdown(srv)

	// select {}
}

// printStartupInfo 打印启动信息
func (a *App) printStartupInfo(conf *config.Config) {
	// 应用信息
	appInfo := map[string]interface{}{
		"应用":  conf.App.Name,
		"环境":  conf.App.Env,
		"端口":  conf.App.Port,
		"数据库": conf.Mysql.Database,
	}
	PrintAligned(appInfo, []string{"应用", "环境", "端口", "数据库"})

	// 服务地址
	port := pkg.IntToString(conf.App.Port)
	network := pkg.Sprintf("%s Address:", flag.NetworkIco)
	swagger := pkg.Sprintf("%s Swagger:", flag.PointerIco)
	testApi := pkg.Sprintf("%s Test API:", flag.PointerIco)
	serverInfo := map[string]interface{}{
		network: pkg.Sprintf("http://0.0.0.0:%s", port),
		swagger: pkg.Sprintf("http://127.0.0.1:%s/swagger/index.html", port),
		testApi: pkg.Sprintf("http://127.0.0.1:%s/ping", port),
	}
	PrintAligned(serverInfo, []string{
		network,
		swagger,
		testApi,
	})

	flag.Successf("Gin server started successfully!")
}

// startHTTPServer 启动http服务
func (a *App) startHttpServer(conf *config.Config) *http.Server {
	port := pkg.IntToString(conf.App.Port)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           a.Engine,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			flag.Errorf("服务启动失败: %s", err.Error())
			os.Exit(1)
		}
	}()

	return srv
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	color.Yellow("服务正在关闭...")

	defer func() {
		app := foundation.GetApp()
		err := app.Stop()
		if err != nil {
			flag.Errorf("关闭应用失败: %s", err.Error())
			return
		}
	}()

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭http服务
	if err := srv.Shutdown(ctx); err != nil {
		flag.Errorf("服务关闭异常: %v", err)
	}
}

// PrintAligned 打印冒号对齐,支持中文
func PrintAligned(data map[string]interface{}, order []string) {
	// 找出最长key的显示宽度
	maxLen := 0
	for k := range data {
		w := runewidth.StringWidth(k)
		if w > maxLen {
			maxLen = w
		}
	}

	for _, k := range order {
		key := ensureEmojiSpace(strings.TrimSuffix(k, ":"))
		padding := maxLen - runewidth.StringWidth(key) + 2
		fmt.Printf("%s:%s%v\n", key, spaces(padding), data[k])
	}
}

// 格式化符合对齐
func ensureEmojiSpace(s string) string {
	r := []rune(s)
	if len(r) > 0 && (r[0] > 0x1F000 && r[0] < 0x1FAFF) {
		if len(r) > 1 && r[1] != ' ' {
			return string(r[0]) + " " + string(r[1:])
		}
	}
	return s
}

// spaces 生成n个空格
func spaces(n int) string {
	if n <= 0 {
		return ""
	}

	return fmt.Sprintf("%*s", n, "")
}

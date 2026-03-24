package bootstrap

import (
	"context"
	"errors"
	"fmt"
	_ "gin/app/listener"
	"gin/app/middleware"
	_ "gin/app/queue/kafka/consumer"
	_ "gin/app/queue/rabbitmq/consumer"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/container"
	"gin/pkg/debugger"
	"gin/pkg/lang"
	"gin/pkg/message"
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

type App struct {
	Engine    *gin.Engine
	Container *container.Container
}

func NewApp(c *container.Container) *App {
	var (
		r    = gin.New()
		conf = config.NewConfig()
	)

	// 运行环境模式 debug模式, test测试模式, release生产模式, 默认是debug,根据当前配置文件读取
	gin.SetMode(conf.App.Mode)

	if conf.App.Env != "production" {
		// 非生产环境允许所有代理
		_ = r.SetTrustedProxies(nil)
	}

	// 设置http请求处理文件上传时可以使用的最大内存为90MB
	r.MaxMultipartMemory = 90 << 20

	// 加载路由
	router.LoadRouters(r)

	return &App{
		Engine:    r,
		Container: c,
	}
}

func (a *App) Run() {
	var conf = config.NewConfig()
	// 设置50：更积极回收
	debug.SetGCPercent(50)
	// 加载翻译
	lang.LoadLang()
	// debugger订阅
	dbg := debugger.NewDebugger(message.GetEventBus())
	dbg.Start()
	defer dbg.Stop()
	middleware.InitRateLimit()

	data := map[string]interface{}{
		"应用":  conf.App.Name,
		"环境":  config.GetString("app.env"),
		"端口":  conf.App.Port,
		"数据库": conf.Mysql.Database,
	}

	// 启动提示
	PrintAligned(data, []string{"应用", "环境", "端口", "数据库"})

	var port = pkg.IntToString(conf.App.Port)
	run := map[string]interface{}{
		flag.Network + " Address:":  "http://0.0.0.0:" + port,
		flag.Pointer + " Swagger:":  "http://127.0.0.1:" + port + "/swagger/index.html",
		flag.Pointer + " Test API:": "http://127.0.0.1:" + port + "/ping",
	}
	PrintAligned(run, []string{flag.Network + " Address:", flag.Pointer + " Swagger:", flag.Pointer + " Test API:"})
	fmt.Println("Gin server started successfully!")

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           a.Engine,
		ReadTimeout:       10 * time.Second, // 设置读取超时
		WriteTimeout:      10 * time.Second, // 设置写入超时
		IdleTimeout:       30 * time.Second, // 设置空闲超时
		ReadHeaderTimeout: 5 * time.Second,  // 设置读取头超时
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			color.Red("服务启动失败: %s", err.Error())
			os.Exit(1)
		}
	}()

	// 等待中断信号以优雅地关闭服务器(设置5秒的超时时间)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	color.Yellow("服务正在关闭...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		color.Red(flag.Error+" 服务关闭异常: %v", err)
	}

	a.shutdown()
	// select {}
}

// 关闭资源
func (a *App) shutdown() {
	middleware.ShutdownRateLimit()

	// ...
	// kafka.Close()
	// rabbitmq.Close()
	// db.Close()
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

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
	"net"
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
	// 获取IP地址
	networkIP := getNetworkIP()
	port := pkg.IntToString(conf.App.Port)

	// 构建数据
	data := map[string]interface{}{
		"应用":  conf.App.Name,
		"环境":  conf.App.Env,
		"端口":  color.YellowString(pkg.IntToString[int64](conf.App.Port)),
		"数据库": conf.Databases.Default,
	}

	// 地址信息
	data[pkg.Sprintf("%s Local Address:", flag.NetworkIco)] = pkg.Sprintf("http://127.0.0.1:%s", port)
	data[pkg.Sprintf("%s Local Swagger:", flag.PointerIco)] = pkg.Sprintf("http://127.0.0.1:%s/swagger/index.html", port)
	data[pkg.Sprintf("%s Local Test API:", flag.PointerIco)] = pkg.Sprintf("http://127.0.0.1:%s/ping", port)

	// 网络地址信息
	if networkIP != "" {
		data[pkg.Sprintf("%s Network Address:", flag.NetworkIco)] = pkg.Sprintf("http://%s:%s", networkIP, port)
		data[pkg.Sprintf("%s Network Swagger:", flag.PointerIco)] = pkg.Sprintf("http://%s:%s/swagger/index.html", networkIP, port)
		data[pkg.Sprintf("%s Network Test API:", flag.PointerIco)] = pkg.Sprintf("http://%s:%s/ping", networkIP, port)
	} else {
		data[pkg.Sprintf("%s Network Address:", flag.NetworkIco)] = "未获取到网络地址"
		data[pkg.Sprintf("%s Network Swagger:", flag.PointerIco)] = "未获取到网络地址"
		data[pkg.Sprintf("%s Network Test API:", flag.PointerIco)] = "未获取到网络地址"
	}

	// 定义显示顺序
	order := []string{
		"应用", "环境", "端口", "数据库",
		pkg.Sprintf("%s Local Address:", flag.NetworkIco),
		pkg.Sprintf("%s Network Address:", flag.NetworkIco),
		pkg.Sprintf("%s Local Swagger:", flag.PointerIco),
		pkg.Sprintf("%s Network Swagger:", flag.PointerIco),
		pkg.Sprintf("%s Local Test API:", flag.PointerIco),
		pkg.Sprintf("%s Network Test API:", flag.PointerIco),
	}

	PrintAligned(data, order)

	flag.Successf("Gin server started successfully!")
}

// PrintAligned 打印冒号对齐,支持中文
func PrintAligned(data map[string]interface{}, order []string) {
	if len(order) == 0 {
		return
	}

	// 找出最长key的显示宽度
	maxLen := 0
	for _, k := range order {
		if val, ok := data[k]; ok && val != nil {
			w := runewidth.StringWidth(k)
			if w > maxLen {
				maxLen = w
			}
		}
	}

	// 打印每一行
	for _, k := range order {
		if val, ok := data[k]; ok && val != nil {
			key := ensureEmojiSpace(strings.TrimSuffix(k, ":"))
			padding := maxLen - runewidth.StringWidth(key) + 2

			// 根据值的类型处理颜色
			switch v := val.(type) {
			case string:
				if strings.HasPrefix(v, "http") {
					// URL使用青色
					fmt.Printf("%s:%s%s\n", key, spaces(padding), color.CyanString(v))
				} else if strings.Contains(k, "Network") && v == "未获取到网络地址" {
					// 未获取到网络地址使用黄色
					fmt.Printf("%s:%s%s\n", key, spaces(padding), color.YellowString(v))
				} else {
					// 普通字符串
					fmt.Printf("%s:%s%s\n", key, spaces(padding), color.YellowString(v))
				}
			default:
				fmt.Printf("%s:%s%v\n", key, spaces(padding), v)
			}
		}
	}
}

// getNetworkIP 获取局域网IP地址
func getNetworkIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			// 排除Docker和虚拟网卡
			if !isVirtualIP(ipnet.IP.String()) {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// isVirtualIP 判断是否为虚拟网卡IP
func isVirtualIP(ip string) bool {
	// 常见的虚拟网卡IP段
	virtualPrefixes := []string{
		"172.17.", // Docker
		"172.18.",
		"172.19.",
		"192.168.99.", // Docker Machine
		"10.0.",       // VPN
	}

	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

// startHttpServer 启动http服务
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

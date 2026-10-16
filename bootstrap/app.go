package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"gin/app/facade"
	"gin/app/provider"
	"gin/app/service"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/errcode"
	"gin/pkg/route"
	"gin/pkg/serviceprovider"
	"gin/router"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"gin/pkg/convert"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-runewidth"
)

// App 应用结构
type App struct {
	Engine      *gin.Engine
	application *serviceprovider.Application
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{
		application: serviceprovider.NewApp(provider.All()...),
	}
}

// Boot 启动应用
func (a *App) Boot() error {
	if a == nil || a.application == nil {
		return errors.New("application not initialized")
	}
	if a.Engine != nil {
		return nil
	}

	if err := a.application.Boot(); err != nil {
		return err
	}

	a.Engine = setupEngine()
	return nil
}

// setupEngine 配置Gin引擎
func setupEngine() *gin.Engine {
	r := gin.New()

	conf := facade.Config()
	errcode.SetLogger(facade.Log())

	// 设置运行模式
	gin.SetMode(conf.App.Mode)

	// 非生产环境允许所有代理
	if conf.App.Env != "production" {
		_ = r.SetTrustedProxies(nil)
	}

	// 设置文件上传最大内存
	r.MaxMultipartMemory = 90 << 20

	// 加载路由
	router.NewRouters(r)

	// 同步路由权限到数据库
	syncRoutePermissions()

	return r
}

// syncRoutePermissions 同步路由权限到数据库
func syncRoutePermissions() {
	permissionKeys := route.PermissionKeys(router.All()...)
	if len(permissionKeys) == 0 {
		return
	}

	svc := service.PermissionService{}
	if err := svc.SyncRoutePermissions(context.Background(), permissionKeys); err != nil {
		flag.Warningf("同步路由权限失败: %s", err.Error())
	}
}

func (a *App) Run() {
	if a == nil || a.Engine == nil {
		flag.Errorf("应用未启动,无法运行服务")
		os.Exit(1)
	}

	// 系统优化设置50：更积极回收
	debug.SetGCPercent(50)

	conf := facade.Config()
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

// SetupCLI 初始化CLI运行环境,不启动后台任务
func SetupCLI() error {
	config.NewConfig()
	flag.SetSilent(true)
	defer flag.SetSilent(false)

	app := serviceprovider.NewApp(provider.All()...)
	if err := app.RegisterProviders(); err != nil {
		return err
	}

	return nil
}

// printStartupInfo 打印启动信息
func (a *App) printStartupInfo(conf *config.Config) {
	// 获取IP地址
	networkIP := getNetworkIP()
	port := convert.IntToString(conf.App.Port)
	localHost := "127.0.0.1"
	networkAddress := "未获取到网络地址"
	networkSwagger := "未获取到网络地址"
	networkTestAPI := "未获取到网络地址"
	if networkIP != "" {
		networkHost := "http://" + networkIP + ":" + port
		networkAddress = networkHost
		networkSwagger = networkHost + "/swagger/index.html"
		networkTestAPI = networkHost + "/ping"
	}

	// 构建启动信息
	items := []startupItem{
		{label: "应用", value: conf.App.Name},
		{label: "环境", value: conf.App.Env},
		{label: "端口", value: color.YellowString(convert.IntToString[int64](conf.App.Port))},
		{label: "数据库", value: conf.Databases.Driver},
		{label: startupLabel(flag.NetworkIco, "Local Address"), value: "http://" + localHost + ":" + port},
		{label: startupLabel(flag.NetworkIco, "Network Address"), value: networkAddress},
		{label: startupLabel(flag.PointerIco, "Local Swagger"), value: "http://" + localHost + ":" + port + "/swagger/index.html"},
		{label: startupLabel(flag.PointerIco, "Network Swagger"), value: networkSwagger},
		{label: startupLabel(flag.PointerIco, "Local Test API"), value: "http://" + localHost + ":" + port + "/ping"},
		{label: startupLabel(flag.PointerIco, "Network Test API"), value: networkTestAPI},
	}

	printAligned(items)

	flag.Successf("Gin server started successfully!")
}

// startupItem 启动信息项
type startupItem struct {
	label string
	value any
}

// printAligned 打印冒号对齐,支持中文
func printAligned(items []startupItem) {
	// 找出最长key的显示宽度
	maxLen := 0
	for _, item := range items {
		if item.value != nil {
			w := runewidth.StringWidth(item.label)
			if w > maxLen {
				maxLen = w
			}
		}
	}

	// 打印每一行
	for _, item := range items {
		if item.value != nil {
			key := ensureEmojiSpace(strings.TrimSuffix(item.label, ":"))
			padding := maxLen - runewidth.StringWidth(key) + 2

			// 根据值的类型处理颜色
			switch v := item.value.(type) {
			case string:
				if strings.HasPrefix(v, "http") {
					// URL使用青色
					fmt.Printf("%s:%s%s\n", key, spaces(padding), color.CyanString(v))
				} else if strings.Contains(item.label, "Network") && v == "未获取到网络地址" {
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

// startupLabel 构建启动信息标签
func startupLabel(icon, label string) string {
	return icon + " " + label + ":"
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
	port := convert.IntToString(conf.App.Port)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           a.Engine,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      0 * time.Second,
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
		err := a.application.Stop()
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

// 格式化符号对齐
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

	return strings.Repeat(" ", n)
}

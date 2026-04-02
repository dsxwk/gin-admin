package route

import (
	"fmt"
	"gin/common/base"
	"gin/pkg"
	"gin/pkg/cli"
	"gin/router"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)

type RouteList struct{}

func (s *RouteList) Name() string {
	return "route:list"
}

func (s *RouteList) Description() string {
	return "路由列表"
}

func (s *RouteList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *RouteList) Execute(args []string) {
	// 初始化Gin引擎(不要Run)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 加载项目路由
	router.LoadRouters(engine)

	// 获取所有路由
	routes := engine.Routes()
	if len(routes) == 0 {
		color.Yellow("暂无注册的路由")
		return
	}

	// 按Path排序
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})

	// 打印路由列表
	fmt.Println(pkg.Sprintf("%-8s %-35s %-40s", "Method", "Path", "Handler"))
	for _, route := range routes {
		// Method颜色: GET=绿色 POST=黄色 PUT=蓝色 DELETE=红色
		var methodColor *color.Color
		switch route.Method {
		case "GET":
			methodColor = color.New(color.FgGreen)
		case "POST":
			methodColor = color.New(color.FgYellow)
		case "PUT":
			methodColor = color.New(color.FgBlue)
		case "DELETE":
			methodColor = color.New(color.FgRed)
		default:
			methodColor = color.New(color.FgWhite)
		}

		// Path颜色: 青色
		pathColor := color.New(color.FgCyan)

		// Handler颜色: 白色(默认)
		handlerColor := color.New(color.FgWhite)

		str := pkg.Sprintf("%s %s %s",
			methodColor.Sprintf("%-8s", route.Method),
			pathColor.Sprintf("%-35s", route.Path),
			handlerColor.Sprintf("%-40s", s.formatHandlerName(route.Handler)),
		)
		fmt.Println(str)
	}

	color.Cyan("总计 %d 条路由\n", len(routes))
}

func init() {
	cli.Register(&RouteList{})
}

func (s *RouteList) formatHandlerName(handler string) string {
	// 去掉 -fm 结尾
	handler = strings.TrimSuffix(handler, "-fm")

	// 去掉 .func1
	return strings.TrimSuffix(handler, ".func1")
}

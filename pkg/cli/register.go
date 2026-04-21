package cli

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"os"
	"sort"
	"strings"
)

var (
	commands = make(map[string]base.Command)
)

func Register(cmd base.Command) {
	name := cmd.Name()
	if _, exists := commands[name]; exists {
		flag.Errorf("Command \"%s\" already registered, skipped.", name)
		os.Exit(1)
	}
	commands[name] = cmd
}

func Get(name string) (base.Command, bool) {
	cmd, exists := commands[name]
	return cmd, exists
}

func Execute() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage("txt")
		return
	}

	// 全局选项
	switch args[0] {
	case "-h", "--help":
		printUsage("txt")
		return
	case "-v", "--version":
		fmt.Println("Gin Cli", getVersion())
		return
	}

	if strings.HasPrefix(args[0], "-f") || strings.HasPrefix(args[0], "--format") {
		format := "txt"

		// 支持三种写法：
		//   -f json
		//   -f=json
		//   --format=json
		if len(args) > 1 && !strings.Contains(args[0], "=") {
			format = args[1]
		} else if strings.Contains(args[0], "=") {
			parts := strings.SplitN(args[0], "=", 2)
			if len(parts) == 2 {
				format = parts[1]
			}
		}
		printUsage(format)
		return
	}

	// 子命令名
	name := args[0]
	cmdArgs := args[1:]

	cmd, exists := Get(name)
	if !exists {
		flag.Errorf("Command \"%s\" is not defined.", name)
		printUsage("txt")
		os.Exit(1)
	}

	// --help 自动打印命令帮助
	for _, arg := range cmdArgs {
		if arg == "-h" || arg == "--help" {
			printCommandHelp(cmd)
			return
		}
	}

	// 交给命令自己解析参数
	cmd.Execute(cmdArgs)
}

// 打印命令列表
func printUsage(format string) {
	switch format {
	case "json":
		printJSON()
	default:
		printText()
	}
}

// 打印文本格式
func printText() {
	fmt.Println("Gin Cli", getVersion())
	fmt.Println()

	color.Yellow("Usage:")
	fmt.Println("  cli [command] [options]")
	fmt.Println()
	color.Yellow("Available commands:")

	// 按前缀分组
	groups := make(map[string][]string)
	maxCmdLen := 0

	// 收集命令并分组
	for name := range commands {
		parts := strings.SplitN(name, ":", 2)
		group := parts[0]
		if len(parts) == 1 {
			group = "other"
		}

		groups[group] = append(groups[group], name)

		if len(name) > maxCmdLen {
			maxCmdLen = len(name)
		}
	}

	// 获取所有组名并排序
	groupNames := make([]string, 0, len(groups))
	for group := range groups {
		if group != "other" {
			groupNames = append(groupNames, group)
		}
	}
	sort.Strings(groupNames)

	// 将other组放在最后
	if _, ok := groups["other"]; ok {
		groupNames = append(groupNames, "other")
	}

	// 计算最大宽度
	maxWidth := 0
	for name := range commands {
		w := runewidth.StringWidth(name)
		if w > maxWidth {
			maxWidth = w
		}
	}

	// 打印Options
	options := [][]string{
		{"-f, --format", "The output format (txt, json) [default: txt]"},
		{"-h, --help", "Display help for the given command"},
		{"-v, --version", "Display cli version"},
	}

	// 计算最大宽度
	optMax := 0
	for _, opt := range options {
		w := runewidth.StringWidth(opt[0])
		if w > optMax {
			optMax = w
		}
	}

	// 打印命令
	for _, group := range groupNames {
		names := groups[group]
		sort.Strings(names)

		color.Yellow("%s:", group)

		var items [][]string
		for _, name := range names {
			cmd := commands[name]
			items = append(items, []string{name, cmd.Description()})
		}

		printAlign(items, 2, lo.Ternary(true, maxWidth, optMax))
	}

	color.Yellow("\nOptions:")
	printAlign(options, 2, lo.Ternary(true, maxWidth, optMax))
}

// 打印单个命令帮助
func printCommandHelp(cmd base.Command) {
	fmt.Println("Gin Cli", getVersion())
	color.Yellow("\nUsage:")
	fmt.Println("  cli [command] [options]")
	color.Yellow("\nCommand:")
	fmt.Printf("  %s  %s\n\n", color.GreenString(cmd.Name()), cmd.Description())

	options := cmd.Help()
	if len(options) == 0 {
		fmt.Println("该命令暂无选项")
		return
	}

	color.Yellow("Options:")
	base.PrintArgs(options)
}

// 打印json格式
func printJSON() {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	var list []map[string]string
	for _, name := range names {
		cmd := commands[name]
		list = append(list, map[string]string{
			"name":        name,
			"description": cmd.Description(),
		})
	}

	data := map[string]interface{}{
		"version":  "Gin Cli " + facade.Config.Get().App.CliVersion,
		"commands": list,
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	color.Green(string(jsonData))
}

func getVersion() string {
	return color.GreenString(facade.Config.Get().App.CliVersion)
}

// 打印对齐
func printAlign(items [][]string, indent int, maxWidth int) {
	for _, item := range items {
		key := item[0]
		val := item[1]

		w := runewidth.StringWidth(key)
		padding := maxWidth - w

		fmt.Printf("%s%s%s  %s\n",
			strings.Repeat(" ", indent),
			color.GreenString(key),
			strings.Repeat(" ", padding),
			val,
		)
	}
}

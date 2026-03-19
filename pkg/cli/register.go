package cli

import (
	"fmt"
	"gin/common/base"
	"gin/common/flag"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
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
		color.Yellow(flag.Warning+"  Command \"%s\" already registered, skipped.", name)
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
		color.Green("Gin CLI v1.0.0")
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
		color.Red(flag.Error+"  Command \"%s\" is not defined.", name)
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
	color.Cyan("Usage: cli [command] [options]")
	fmt.Println()
	color.Yellow("Available commands:")

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := commands[name]
		fmt.Printf("  %s%s\n", color.GreenString(fmt.Sprintf("%-25s", name)), cmd.Description())
	}

	fmt.Println()
	color.Yellow("Options:")
	fmt.Println(color.GreenString("  -f, --format        ") + "The output format (txt, json) [default: txt]")
	fmt.Println(color.GreenString("  -h, --help          ") + "Display help for the given command")
	fmt.Println(color.GreenString("  -v, --version       ") + "Display CLI version")
}

// 打印单个命令帮助
func printCommandHelp(cmd base.Command) {
	fmt.Printf("\n%s - %s\n\n", color.GreenString(cmd.Name()), cmd.Description())

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
		"version":  "Gin CLI v2.0.0",
		"commands": list,
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	color.Green(string(jsonData))
}

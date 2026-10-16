package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin/app/facade"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
)

// Registry 命令注册表
type Registry struct {
	commands map[string]base.Handler
}

// NewRegistry 创建命令注册表
func NewRegistry(items ...base.Handler) *Registry {
	registry := &Registry{
		commands: make(map[string]base.Handler, len(items)),
	}
	for _, item := range items {
		registry.Register(item)
	}
	return registry
}

// flagParser flag解析器
type flagParser interface {
	ParseFlags(name string, args []string, opts []base.CommandOption) (map[string]string, error)
}

// Register 注册命令
func (r *Registry) Register(cmd base.Handler) {
	if r == nil || cmd == nil {
		return
	}

	name := cmd.Name()
	if _, exists := r.commands[name]; exists {
		terminal.Errorf("Command \"%s\" already registered, skipped.", name)
		os.Exit(1)
	}
	r.commands[name] = cmd
}

// Get 获取命令
func (r *Registry) Get(name string) (base.Handler, bool) {
	if r == nil {
		return nil, false
	}
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute 创建注册表并执行命令
func Execute(items ...base.Handler) {
	NewRegistry(items...).Execute(os.Args[1:])
}

// Execute 执行命令
func (r *Registry) Execute(args []string) {
	if len(args) == 0 {
		r.printUsage("txt")
		return
	}

	// 全局选项
	switch args[0] {
	case "-h", "--help":
		r.printUsage("txt")
		return
	case "-v", "--version":
		printVersion()
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
		r.printUsage(format)
		return
	}

	// 子命令名
	name := args[0]
	cmdArgs := args[1:]

	cmd, exists := r.Get(name)
	if !exists {
		terminal.Errorf("Command \"%s\" is not defined.", name)
		r.printUsage("txt")
		os.Exit(1)
	}

	// --help 自动打印命令帮助
	for _, arg := range cmdArgs {
		if arg == "-h" || arg == "--help" {
			printCommandHelp(cmd)
			return
		}
	}

	// 统一解析flag后执行命令
	if fp, ok := cmd.(flagParser); ok {
		values, err := fp.ParseFlags(cmd.Name(), cmdArgs, cmd.Help())
		if err != nil {
			terminal.Errorf("%s", err.Error())
			base.PrintArgs(cmd.Help())
			os.Exit(1)
		}
		if err = cmd.Handle(values); err != nil {
			terminal.Errorf("%s", err.Error())
			os.Exit(1)
		}
	} else {
		if err := cmd.Handle(nil); err != nil {
			terminal.Errorf("%s", err.Error())
			os.Exit(1)
		}
	}
}

// printVersion 打印版本信息
func printVersion() {
	const asciiDsx = `
  ██████  ██████ ██   ██                                
  ██   ██ ██      ██ ██                                
  ██   ██ ██████   ███                                 
  ██   ██     ██  ██ ██                                 
  ██████  ██████ ██   ██ 
  
Gin Cli %s, built with Go %s`

	fmt.Printf(asciiDsx, getVersion(), runtime.Version())
	fmt.Println()
}

// 打印命令列表
func (r *Registry) printUsage(format string) {
	switch format {
	case "json":
		r.printJSON()
	default:
		r.printText()
	}
}

// 打印文本格式
func (r *Registry) printText() {
	printVersion()
	fmt.Println()

	color.Yellow("Usage:")
	fmt.Println("  cli [command] [options]")
	fmt.Println()
	color.Yellow("Available commands:")

	// 按前缀分组
	groups := make(map[string][]string)
	maxCmdLen := 0

	// 收集命令并分组
	for name := range r.commands {
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
	for name := range r.commands {
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
			cmd := r.commands[name]
			items = append(items, []string{name, cmd.Description()})
		}

		printAlign(items, 2, lo.Ternary(maxWidth > optMax, maxWidth, optMax))
	}

	color.Yellow("\nOptions:")
	printAlign(options, 2, lo.Ternary(maxWidth > optMax, maxWidth, optMax))
}

// 打印单个命令帮助
func printCommandHelp(cmd base.Handler) {
	printVersion()
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
func (r *Registry) printJSON() {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	sort.Strings(names)

	var list []map[string]string
	for _, name := range names {
		cmd := r.commands[name]
		list = append(list, map[string]string{
			"name":        name,
			"description": cmd.Description(),
		})
	}

	data := map[string]any{
		"version":  "Gin Cli " + facade.Config().App.CliVersion,
		"commands": list,
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	color.Green(string(jsonData))
}

func getVersion() string {
	return color.GreenString(facade.Config().App.CliVersion)
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

// AddItem 添加注册项到指定列表
func AddItem(file, markerText, itemText string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	item := []byte(itemText)
	if bytes.Contains(content, item) {
		return nil
	}

	marker := []byte(markerText)
	start := bytes.Index(content, marker)
	if start < 0 {
		return fmt.Errorf("未找到注册列表")
	}

	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	bodyStart := start + len(marker)
	closeMarker := []byte(newline + "\t}")
	end := bytes.Index(content[bodyStart:], closeMarker)
	if end < 0 {
		return fmt.Errorf("注册列表格式不正确")
	}

	insertAt := bodyStart + end
	itemLine := []byte(newline + "\t\t" + string(item) + ",")
	result := make([]byte, 0, len(content)+len(itemLine))
	result = append(result, content[:insertAt]...)
	result = append(result, itemLine...)
	result = append(result, content[insertAt:]...)

	return os.WriteFile(file, result, 0644)
}

// AddImport 添加注册文件导入
func AddImport(file, importPath string) (string, error) {
	alias := importAlias(importPath)
	if err := addImportLine(file, importPath, alias); err != nil {
		return "", err
	}

	return alias, nil
}

// addImportLine 添加导入行
func addImportLine(file, importPath, alias string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	quotedPath := fmt.Sprintf("%q", importPath)
	if bytes.Contains(content, []byte(quotedPath)) {
		_, err = ensureImportAlias(file, content, quotedPath, alias)
		return err
	}

	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	marker := []byte("import (")
	start := bytes.Index(content, marker)
	if start < 0 {
		return fmt.Errorf("未找到导入列表")
	}

	bodyStart := start + len(marker)
	closeMarker := []byte(newline + ")")
	end := bytes.Index(content[bodyStart:], closeMarker)
	if end < 0 {
		return fmt.Errorf("导入列表格式不正确")
	}

	insertAt := bodyStart + end
	importLine := []byte(newline + "\t" + alias + " " + quotedPath)
	result := make([]byte, 0, len(content)+len(importLine))
	result = append(result, content[:insertAt]...)
	result = append(result, importLine...)
	result = append(result, content[insertAt:]...)

	return os.WriteFile(file, result, 0644)
}

// ensureImportAlias 确保导入使用别名
func ensureImportAlias(file string, content []byte, quotedPath, alias string) (string, error) {
	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	lines := strings.Split(string(content), newline)
	for index, line := range lines {
		if !strings.Contains(line, quotedPath) {
			continue
		}

		current := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), quotedPath))
		if current != "" {
			return current, nil
		}

		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		lines[index] = indent + alias + " " + quotedPath
		if err := os.WriteFile(file, []byte(strings.Join(lines, newline)), 0644); err != nil {
			return "", err
		}

		return alias, nil
	}

	return alias, nil
}

// importAlias 生成注册导入别名
func importAlias(importPath string) string {
	alias := lo.CamelCase(path.Base(importPath))
	if alias == "" {
		return "dynamic"
	}

	return alias
}

// ImportPath 生成注册导入路径
func ImportPath(file string) string {
	dir := filepath.ToSlash(filepath.Dir(file))
	dir = strings.TrimPrefix(dir, "./")
	return "gin/" + strings.TrimPrefix(dir, "/")
}

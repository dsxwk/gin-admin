package base

import (
	"bufio"
	"context"
	"fmt"
	"gin/pkg/pathx"
	"gin/pkg/stringx"
	"gin/pkg/terminal"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
)

// Config 命令配置
type Config struct {
	Name        string
	Description string
	Options     []CommandOption
	Cron        string
	CronValues  map[string]string
}

// cronOption 通用cron参数
var cronOption = CommandOption{
	Flag: Flag{
		Long: "cron",
	},
	Desc: "Cron表达式,如: */5 * * * * *,为空时手动执行",
}

var commandCronParser = cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// Command 基础命令
type Command struct {
	config Config
}

// NewCommand 创建基础命令
func NewCommand(config Config) Command {
	if config.Description == "" {
		config.Description = config.Name
	}
	if !hasOption(config.Options, cronOption.Flag.Long) {
		config.Options = append(config.Options, cronOption)
	}
	return Command{config: config}
}

// Name 获取命令名称
func (b *Command) Name() string {
	return b.config.Name
}

// Description 获取命令描述
func (b *Command) Description() string {
	return b.config.Description
}

// Flag 定义短长参数名
type Flag struct {
	Short   string
	Long    string
	Default string
}

// CommandOption 用于命令选项描述
type CommandOption struct {
	Flag     Flag   // flag定义
	Desc     string // 描述
	Required bool   // 是否必填
}

// Handler 命令处理器
type Handler interface {
	Name() string                          // 命令名称，如 "make:controller"
	Description() string                   // 命令描述
	Handle(values map[string]string) error // 执行逻辑
	Help() []CommandOption                 // 获取命令帮助信息
}

// Scheduled 定时命令接口
type Scheduled interface {
	Cron() string
	CronValues() map[string]string
}

// Help 获取命令帮助
func (b *Command) Help() []CommandOption {
	return b.config.Options
}

// Cron 获取自动调度表达式
func (b *Command) Cron() string {
	return b.config.Cron
}

// CronValues 获取自动调度参数
func (b *Command) CronValues() map[string]string {
	return cloneValues(b.config.CronValues)
}

// CronManager 控制台定时任务管理器
type CronManager struct {
	scheduler *cron.Cron
	count     int
}

// NewCronManager 创建控制台定时任务管理器
func NewCronManager(commands ...Handler) (*CronManager, error) {
	scheduler := cron.New(cron.WithParser(commandCronParser))
	count := 0

	for _, item := range commands {
		scheduled, ok := item.(Scheduled)
		if !ok {
			continue
		}

		expression := strings.TrimSpace(scheduled.Cron())
		if expression == "" {
			continue
		}
		if _, err := commandCronParser.Parse(expression); err != nil {
			return nil, fmt.Errorf("命令%s的cron表达式无效:%w", item.Name(), err)
		}

		handler := item
		values := scheduled.CronValues()
		if _, err := scheduler.AddFunc(expression, func() {
			if err := handler.Handle(cloneValues(values)); err != nil {
				terminal.Errorf("命令%s定时执行失败:%s", handler.Name(), err.Error())
			}
		}); err != nil {
			return nil, fmt.Errorf("注册命令%s的cron任务失败:%w", item.Name(), err)
		}
		count++
	}

	return &CronManager{
		scheduler: scheduler,
		count:     count,
	}, nil
}

// Run 启动控制台定时任务
func (m *CronManager) Run(ctx context.Context) error {
	if m == nil || m.scheduler == nil || m.count == 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	m.scheduler.Start()
	terminal.Infof("控制台定时任务已启动,任务数:%d", m.count)
	<-ctx.Done()
	return nil
}

// Stop 停止控制台定时任务
func (m *CronManager) Stop() error {
	if m == nil || m.scheduler == nil {
		return nil
	}
	m.scheduler.Stop()
	terminal.Infof("控制台定时任务已停止")
	return nil
}

// Name 获取定时任务管理器名称
func (m *CronManager) Name() string {
	return "console_cron"
}

// Len 获取定时任务数量
func (m *CronManager) Len() int {
	if m == nil {
		return 0
	}
	return m.count
}

// ExecuteWithCron 执行命令或启动cron调度
func ExecuteWithCron(command Handler, values map[string]string) error {
	if command == nil {
		return fmt.Errorf("命令未初始化")
	}

	expression := strings.TrimSpace(values["cron"])
	if expression == "" {
		return command.Handle(values)
	}
	if _, err := commandCronParser.Parse(expression); err != nil {
		return fmt.Errorf("解析cron表达式失败:%w", err)
	}

	scheduler := cron.New(cron.WithParser(commandCronParser))
	if _, err := scheduler.AddFunc(expression, func() {
		if err := command.Handle(cloneValues(values)); err != nil {
			terminal.Errorf("命令%s定时执行失败:%s", command.Name(), err.Error())
		}
	}); err != nil {
		return fmt.Errorf("注册cron任务失败:%w", err)
	}

	scheduler.Start()
	terminal.Infof("命令%s已启动定时执行,cron:%s", command.Name(), expression)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	scheduler.Stop()
	terminal.Infof("命令%s定时执行已停止", command.Name())
	return nil
}

// cloneValues 复制命令参数
func cloneValues(values map[string]string) map[string]string {
	if len(values) == 0 {
		return values
	}
	result := make(map[string]string, len(values))
	maps.Copy(result, values)
	return result
}

// ParseFlags flag解析
func (b *Command) ParseFlags(name string, args []string, opts []CommandOption) (map[string]string, error) {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.SetOutput(nil)

	flagRefs := make(map[string]*string)
	for _, opt := range opts {
		if opt.Flag.Short == "" {
			flagRefs[opt.Flag.Long] = fs.String(opt.Flag.Long, opt.Flag.Default, opt.Desc)
			continue
		}
		flagRefs[opt.Flag.Long] = fs.StringP(opt.Flag.Long, opt.Flag.Short, opt.Flag.Default, opt.Desc)
	}

	if err := fs.Parse(normalizeArgs(args, opts)); err != nil {
		return nil, fmt.Errorf("参数解析失败:%w", err)
	}

	values := make(map[string]string)
	for key, ref := range flagRefs {
		values[key] = strings.TrimPrefix(*ref, "=")
	}

	for _, opt := range opts {
		if opt.Required && values[opt.Flag.Long] == "" {
			return nil, fmt.Errorf("参数%s不能为空", optionLabel(opt))
		}
	}

	terminal.Successf("执行命令: %s %s", name, b.FormatArgs(values))
	return values, nil
}

// normalizeArgs 规范化单横线长参数
func normalizeArgs(args []string, opts []CommandOption) []string {
	if len(args) == 0 {
		return args
	}

	longFlags := make(map[string]struct{}, len(opts))
	for _, opt := range opts {
		longFlags[opt.Flag.Long] = struct{}{}
	}

	result := make([]string, len(args))
	for index, arg := range args {
		result[index] = arg
		if len(arg) < 3 || !strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			continue
		}

		name := strings.TrimPrefix(arg, "-")
		value := ""
		if equalIndex := strings.Index(name, "="); equalIndex >= 0 {
			value = name[equalIndex:]
			name = name[:equalIndex]
		}
		if _, exists := longFlags[name]; exists {
			result[index] = "--" + name + value
		}
	}

	return result
}

// PrintArgs 打印参数
func PrintArgs(opts []CommandOption) {
	maxFlagWidth := 0
	maxDescWidth := 0
	for _, opt := range opts {
		flagStr := optionLabel(opt)
		if w := runewidth.StringWidth(flagStr); w > maxFlagWidth {
			maxFlagWidth = w
		}
		if w := runewidth.StringWidth(opt.Desc); w > maxDescWidth {
			maxDescWidth = w
		}
	}

	for _, opt := range opts {
		flagStr := optionLabel(opt)
		descStr := opt.Desc
		colFlag := color.GreenString(flagStr)
		colDesc := descStr
		flagPad := maxFlagWidth - runewidth.StringWidth(flagStr) + 2
		descPad := maxDescWidth - runewidth.StringWidth(descStr) + 2

		required := color.GreenString("required:false")
		if opt.Required {
			required = color.RedString("required:true")
		}

		fmt.Printf("  %s%s%s%s%s\n",
			colFlag,
			stringx.Spaces(flagPad),
			colDesc,
			stringx.Spaces(descPad),
			required,
		)
	}
}

// hasOption 判断参数是否存在
func hasOption(options []CommandOption, name string) bool {
	for _, option := range options {
		if option.Flag.Long == name {
			return true
		}
	}
	return false
}

// optionLabel 获取参数显示名称
func optionLabel(option CommandOption) string {
	if option.Flag.Short == "" {
		return "--" + option.Flag.Long
	}
	return fmt.Sprintf("-%s, --%s", option.Flag.Short, option.Flag.Long)
}

// FormatArgs 格式化参数
func (b *Command) FormatArgs(args map[string]string) string {
	keys := make([]string, 0, len(args))
	for key := range args {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		value := args[key]
		if value != "" {
			builder.WriteString(fmt.Sprintf("--%s=%s ", key, value))
		}
	}
	return builder.String()
}

// StringToBool 将字符串安全地转换为布尔值
func (b *Command) StringToBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off", "":
		return false
	default:
		return false
	}
}

// ExitError 输出错误并退出
func (b *Command) ExitError(msg string) {
	terminal.Errorf("%s", msg)
	os.Exit(1)
}

// GetMakeFile 获取make文件
func (b *Command) GetMakeFile(file, kind string) string {
	file, err := cleanRelativePath(file)
	if err != nil {
		b.ExitError(err.Error())
		return ""
	}

	switch kind {
	case "router":
		file = filepath.Join("router", file)
	case "ws":
		file = filepath.Join("websocket", file)
	default:
		file = filepath.Join("app", kind, file)
	}

	return file + ".go"
}

// GetTemplate 获取模版文件
func (b *Command) GetTemplate(kind string) string {
	switch kind {
	case "model", "console", "controller", "service", "request", "middleware", "router", "event", "listener", "facade", "provider", "consumer", "producer", "enum", "errcode", "mcp", "es", "ws", "grpc_request", "grpc_service", "grpc_proto":
		return filepath.Join(pathx.RootPath(), "pkg", "cli", "template", kind+".tpl")
	default:
		b.ExitError("未找到 " + kind + " 模版文件")
		return ""
	}
}

// CheckDirAndFile 检查目录和文件
func (b *Command) CheckDirAndFile(file string) *os.File {
	file, err := cleanRelativePath(file)
	if err != nil {
		terminal.Errorf("%s", err.Error())
		return nil
	}

	dir := filepath.Dir(file)
	if err = os.MkdirAll(dir, 0755); err != nil {
		terminal.Errorf("Failed to create directory: %s", err)
		return nil
	}

	if _, statErr := os.Stat(file); statErr == nil {
		fmt.Printf("%s文件 %s 已存在,是否覆盖?(%s/%s): ",
			color.YellowString(terminal.Warning()),
			color.CyanString(file),
			color.GreenString("Y"),
			color.RedString("N"),
		)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if !b.StringToBool(strings.TrimSpace(input)) {
			fmt.Println("操作已取消")
			return nil
		}
	} else if !os.IsNotExist(statErr) {
		terminal.Errorf("Failed to check file: %s", statErr.Error())
		return nil
	}

	terminal.Successf(terminal.FileIco+" 创建文件: %s", color.CyanString(file))
	f, err := os.Create(file)
	if err != nil {
		terminal.Errorf("Failed to create file: %s", err.Error())
		return nil
	}
	return f
}

// cleanRelativePath 校验并清理相对路径
func cleanRelativePath(file string) (string, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	file = filepath.ToSlash(file)
	if filepath.IsAbs(file) || filepath.VolumeName(file) != "" || strings.HasPrefix(file, "/") {
		return "", fmt.Errorf("文件路径必须为相对路径:%s", file)
	}

	clean := filepath.Clean(filepath.FromSlash(file))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("文件路径不能超出项目目录:%s", file)
	}
	return clean, nil
}

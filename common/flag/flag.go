package flag

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

var (
	printMu sync.Mutex
	printed bool
)

// 预定义颜色函数
var (
	// 错误: 红色背景 + 白色文字(前缀),红色文字(内容)
	errorLabel = color.New(color.BgRed, color.FgWhite, color.Bold)
	errorText  = color.New(color.FgRed)

	// 警告: 黄色背景 + 黑色文字(前缀),黄色文字(内容)
	warningLabel = color.New(color.BgYellow, color.FgBlack, color.Bold)
	warningText  = color.New(color.FgYellow)

	// 成功: 绿色背景 + 白色文字(前缀),绿色文字(内容)
	successLabel = color.New(color.BgGreen, color.FgWhite, color.Bold)
	successText  = color.New(color.FgGreen)

	// 信息: 蓝色背景 + 白色文字(前缀),蓝色文字(内容)
	infoLabel = color.New(color.BgBlue, color.FgWhite, color.Bold)
	infoText  = color.New(color.FgBlue)
)

const (
	ErrorIco   = "❌"
	SuccessIco = "✅"
	WarningIco = "⚠️"
	InfoIco    = "ℹ️"
	LoadingIco = "🔄"
	NetworkIco = "🌐"
	PointerIco = "👉"
	FileIco    = "📄"
)

// Error 错误前缀(红色背景)
// 使用示例: fmt.Printf("%s数据库连接失败: %v", flag.Error(), err)
func Error() string {
	return color.New(color.BgRed, color.FgWhite, color.Bold).Sprint("  ERROR  ") + " "
}

// Warning 警告前缀(黄色背景)
func Warning() string {
	return color.New(color.BgYellow, color.FgBlack, color.Bold).Sprint(" WARNING ") + " "
}

// Success 成功前缀(绿色背景)
func Success() string {
	return color.New(color.BgGreen, color.FgWhite, color.Bold).Sprint(" SUCCESS ") + " "
}

// Info 信息前缀(蓝色背景)
func Info() string {
	return color.New(color.BgBlue, color.FgWhite, color.Bold).Sprint("   INFO  ") + " "
}

// Errorf 错误日志(前缀背景色 + 文字颜色)
// 使用示例: flag.Errorf("%s数据库连接失败: %v", "mysql", err)
func Errorf(format string, args ...any) {
	printLine(errorLabel, errorText, "  ERROR  ", format, args...)
}

// Warningf 警告日志
func Warningf(format string, args ...any) {
	printLine(warningLabel, warningText, " WARNING ", format, args...)
}

// Successf 成功日志
func Successf(format string, args ...any) {
	printLine(successLabel, successText, " SUCCESS ", format, args...)
}

// Infof 信息日志
func Infof(format string, args ...any) {
	printLine(infoLabel, infoText, "   INFO  ", format, args...)
}

// printLine 原子输出日志
func printLine(label, text *color.Color, prefix, format string, args ...any) {
	printMu.Lock()
	defer printMu.Unlock()
	line := label.Sprint(prefix) + text.Sprintf(" "+format, args...) + "\n"
	_, _ = fmt.Fprint(color.Output, line)
}

// ErrorEmoji 带Emoji的错误前缀
func ErrorEmoji() string {
	return fmt.Sprintf("%s %s ", ErrorIco, color.New(color.BgRed, color.FgWhite, color.Bold).Sprint(" ERROR  "))
}

// WarningEmoji 带Emoji的警告前缀
func WarningEmoji() string {
	return fmt.Sprintf("%s %s ", WarningIco, color.New(color.BgYellow, color.FgBlack, color.Bold).Sprint(" WARNING"))
}

// SuccessEmoji 带Emoji的成功前缀
func SuccessEmoji() string {
	return fmt.Sprintf("%s %s ", SuccessIco, color.New(color.BgGreen, color.FgWhite, color.Bold).Sprint(" SUCCESS"))
}

// InfoEmoji 带Emoji的信息前缀
func InfoEmoji() string {
	return fmt.Sprintf("%s %s ", InfoIco, color.New(color.BgBlue, color.FgWhite, color.Bold).Sprint(" INFO   "))
}

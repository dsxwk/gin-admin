package flag

import (
	"fmt"
	"github.com/fatih/color"
)

// 预定义颜色函数
var (
	// 错误: 红色背景 + 白色文字(前缀),红色文字(内容)
	errorBg   = color.New(color.BgRed, color.FgWhite, color.Bold)
	errorText = color.New(color.FgRed)

	// 警告: 黄色背景 + 黑色文字(前缀),黄色文字(内容)
	warningBg   = color.New(color.BgYellow, color.FgBlack, color.Bold)
	warningText = color.New(color.FgYellow)

	// 成功: 绿色背景 + 白色文字(前缀),绿色文字(内容)
	successBg   = color.New(color.BgGreen, color.FgWhite, color.Bold)
	successText = color.New(color.FgGreen)

	// 信息: 蓝色背景 + 白色文字(前缀),蓝色文字(内容)
	infoBg   = color.New(color.BgBlue, color.FgWhite, color.Bold)
	infoText = color.New(color.FgBlue)
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

// Error 返回错误前缀(红色背景)
// 使用示例: fmt.Printf("%s数据库连接失败: %v", flag.Error(), err)
func Error() string {
	return color.New(color.BgRed, color.FgWhite, color.Bold).Sprint("  ERROR  ") + " "
}

// Warning 返回警告前缀(黄色背景)
func Warning() string {
	return color.New(color.BgYellow, color.FgBlack, color.Bold).Sprint(" WARNING ") + " "
}

// Success 返回成功前缀(绿色背景)
func Success() string {
	return color.New(color.BgGreen, color.FgWhite, color.Bold).Sprint(" SUCCESS ") + " "
}

// Info 返回信息前缀(蓝色背景)
func Info() string {
	return color.New(color.BgBlue, color.FgWhite, color.Bold).Sprint("   INFO  ") + " "
}

// Errorf 输出错误日志(前缀背景色 + 文字颜色)
// 使用示例: flag.Errorf("%s数据库连接失败: %v", "mysql", err)
func Errorf(format string, args ...interface{}) {
	_, err := errorBg.Print("  ERROR  ")
	if err != nil {
		return
	}
	_, err = errorText.Printf(" "+format+"\n", args...)
	if err != nil {
		return
	}
}

// Warningf 输出警告日志
func Warningf(format string, args ...interface{}) {
	_, err := warningBg.Print(" WARNING ")
	if err != nil {
		return
	}
	_, err = warningText.Printf(" "+format+"\n", args...)
	if err != nil {
		return
	}
}

// Successf 输出成功日志
func Successf(format string, args ...interface{}) {
	_, err := successBg.Print(" SUCCESS ")
	if err != nil {
		return
	}
	_, err = successText.Printf(" "+format+"\n", args...)
	if err != nil {
		return
	}
}

// Infof 输出信息日志
func Infof(format string, args ...interface{}) {
	_, err := infoBg.Print("   INFO  ")
	if err != nil {
		return
	}
	_, err = infoText.Printf(" "+format+"\n", args...)
	if err != nil {
		return
	}
}

// ErrorEmoji 返回带Emoji的错误前缀
func ErrorEmoji() string {
	return fmt.Sprintf("%s %s ", ErrorIco, color.New(color.BgRed, color.FgWhite, color.Bold).Sprint(" ERROR  "))
}

// WarningEmoji 返回带Emoji的警告前缀
func WarningEmoji() string {
	return fmt.Sprintf("%s %s ", WarningIco, color.New(color.BgYellow, color.FgBlack, color.Bold).Sprint(" WARNING"))
}

// SuccessEmoji 返回带Emoji的成功前缀
func SuccessEmoji() string {
	return fmt.Sprintf("%s %s ", SuccessIco, color.New(color.BgGreen, color.FgWhite, color.Bold).Sprint(" SUCCESS"))
}

// InfoEmoji 返回带Emoji的信息前缀
func InfoEmoji() string {
	return fmt.Sprintf("%s %s ", InfoIco, color.New(color.BgBlue, color.FgWhite, color.Bold).Sprint(" INFO   "))
}

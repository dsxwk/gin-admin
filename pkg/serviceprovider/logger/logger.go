package logger

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const StackTraceKey = "stackTrace"

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
	// 全局日志级别(支持动态修改)
	logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

// Logger 包装器
type Logger struct {
	*zap.Logger
}

func NewLogger(conf *config.Config) *Logger {
	loggerOnce.Do(func() {
		setLogLevel(strings.ToLower(conf.Log.Level))

		// 确保日志目录存在
		logDir := filepath.Join(config.RootPath(), "storage", "logs")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			flag.Errorf("创建日志目录失败: %v", err)
			os.Exit(1)
		}

		// 动态日志路径
		logPath := filepath.Join(logDir, "gin.log")

		// 日志切割
		lumberJackLogger := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxDay,
			Compress:   true,
		}

		// 编码配置
		encoderConfig := zap.NewProductionEncoderConfig()
		// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.CallerKey = "caller"
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		// 堆栈
		encoderConfig.StacktraceKey = StackTraceKey

		// 创建encoder,同时输出到文件 + 控制台
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		consoleEncoderConfig := encoderConfig
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		// 动态设置日志级别
		level := logLevel

		// 创建核心
		fileCore := StructuredStackCore{
			Core: zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), level),
		}
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		core := zapcore.NewTee(
			fileCore,
			consoleCore,
		)

		// 初始化 Logger
		zapLogger := zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel), // 自动为error级别以上日志添加堆栈
		)
		zap.ReplaceGlobals(zapLogger) // 替换全局 zap.L()
		loggerInstance = &Logger{zapLogger}
	})

	return loggerInstance
}

// 设置日志级别
func setLogLevel(level string) {
	parsed, err := zapcore.ParseLevel(level)
	if err != nil {
		parsed = zapcore.InfoLevel
	}
	logLevel.SetLevel(parsed)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level string) {
	setLogLevel(strings.ToLower(level))
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() string {
	return logLevel.Level().String()
}

func (l *Logger) WithDebugger(c context.Context) *zap.Logger {
	if l == nil || l.Logger == nil {
		return zap.NewNop()
	}
	if c == nil {
		c = context.Background()
	}

	var ms float64
	if start, ok := c.Value(ctxkey.StartTimeKey).(time.Time); ok {
		ms = float64(time.Since(start).Milliseconds())
	}
	if v, ok := c.Value(ctxkey.MsKey).(float64); ok {
		ms = v
	}

	traceId := getString(c, ctxkey.TraceIDKey)
	trace, _ := debugger.Store.Get(traceId)

	return l.Logger.With(
		zap.String(ctxkey.TraceIDKey, traceId),
		zap.String(ctxkey.IpKey, getString(c, ctxkey.IpKey)),
		zap.String(ctxkey.PathKey, getString(c, ctxkey.PathKey)),
		zap.String(ctxkey.MethodKey, getString(c, ctxkey.MethodKey)),
		zap.Any(ctxkey.ParamsKey, c.Value(ctxkey.ParamsKey)),
		zap.Float64(ctxkey.MsKey, ms),
		zap.Any(ctxkey.DebuggerKey, trace),
	)
}

// 防止panic
func getString(c context.Context, key string) string {
	if c == nil {
		return "unknown"
	}
	if v, ok := c.Value(key).(string); ok {
		return v
	}
	return "unknown"
}

// StructuredStackCore 将自动堆栈编码为字符串数组
type StructuredStackCore struct {
	zapcore.Core
}

// Check 注册结构化堆栈核心
func (c StructuredStackCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(entry.Level) {
		return checked
	}
	return checked.AddCore(entry, c)
}

// With 添加结构化堆栈核心字段
func (c StructuredStackCore) With(fields []zapcore.Field) zapcore.Core {
	return StructuredStackCore{Core: c.Core.With(fields)}
}

// Write 写入结构化堆栈
func (c StructuredStackCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Stack == "" {
		return c.Core.Write(entry, fields)
	}

	lines := splitStack(entry.Stack)
	entry.Stack = ""
	fields = append(fields, zap.Array(StackTraceKey, zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, line := range lines {
			enc.AppendString(line)
		}
		return nil
	})))

	return c.Core.Write(entry, fields)
}

// splitStack 拆分堆栈行
func splitStack(stack string) []string {
	parts := strings.Split(strings.TrimSpace(stack), "\n")
	lines := make([]string, 0, len(parts))
	for _, line := range parts {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

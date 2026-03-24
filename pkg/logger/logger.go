package logger

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/common/trace"
	"gin/config"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
	conf           = config.NewConfig()
	loggerInstance *Logger
	loggerOnce     sync.Once
	// 全局日志级别(支持动态修改)
	logLevel zap.AtomicLevel
)

// Logger 包装器
type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	loggerOnce.Do(func() {
		// 初始化日志级别(默认info)
		logLevel = zap.NewAtomicLevel()
		setLogLevel(strings.ToLower(conf.Log.Level))

		// 确保日志目录存在
		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			color.Red(flag.Error+"  创建日志目录失败:", err)
			os.Exit(1)
		}

		// 动态日志路径
		logPath := filepath.Join(logDir, time.Now().Format("2006-01")+".log")

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
		encoderConfig.StacktraceKey = "stackTrace"

		// 创建encoder,同时输出到文件 + 控制台
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		// 动态设置日志级别
		level := logLevel                                   // 文件跟随全局
		consoleLevel := zap.NewAtomicLevelAt(zap.InfoLevel) // 控制台默认info

		// 创建核心
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel),
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
	switch level {
	case "debug":
		logLevel.SetLevel(zap.DebugLevel)
	case "info":
		logLevel.SetLevel(zap.InfoLevel)
	case "warn":
		logLevel.SetLevel(zap.WarnLevel)
	case "error":
		logLevel.SetLevel(zap.ErrorLevel)
	case "dpanic":
		logLevel.SetLevel(zap.DPanicLevel)
	case "panic":
		logLevel.SetLevel(zap.PanicLevel)
	case "fatal":
		logLevel.SetLevel(zap.FatalLevel)
	default:
		logLevel.SetLevel(zap.InfoLevel)
	}
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
	var ms float64
	if start, ok := c.Value(ctxkey.StartTimeKey).(time.Time); ok {
		ms = float64(time.Since(start).Milliseconds())
	}
	if v, ok := c.Value(ctxkey.MsKey).(float64); ok {
		ms = v
	}

	traceId := getString(c, ctxkey.TraceIdKey)

	return l.Logger.With(
		zap.String(ctxkey.TraceIdKey, traceId),
		zap.String(ctxkey.IpKey, getString(c, ctxkey.IpKey)),
		zap.String(ctxkey.PathKey, getString(c, ctxkey.PathKey)),
		zap.String(ctxkey.MethodKey, getString(c, ctxkey.MethodKey)),
		zap.Any(ctxkey.ParamsKey, c.Value(ctxkey.ParamsKey)),
		zap.Float64(ctxkey.MsKey, ms),
		zap.Any(ctxkey.DebuggerKey, trace.Store.Get(traceId)),
	)
}

// 防止panic
func getString(c context.Context, key string) string {
	if v, ok := c.Value(key).(string); ok {
		return v
	}
	return ""
}

type StackTrace struct{}

func (s StackTrace) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	stack := strings.Split(string(debug.Stack()), "\n")
	return enc.AddArray("stack", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
		for _, line := range stack {
			line = strings.TrimSpace(line)
			if line != "" {
				arr.AppendString(line)
			}
		}
		return nil
	}))
}

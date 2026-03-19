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
)

// Logger 包装器
type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	loggerOnce.Do(func() {
		// 确保日志目录存在
		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			color.Red(flag.Error+"  创建日志目录失败:", err)
			os.Exit(1)
		}

		// 动态日志路径
		logPath := filepath.Join(logDir, time.Now().Format("2006-01")+".log")

		// 滚动日志配置
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
		// 格式化堆栈输出(多行缩进)
		encoderConfig.StacktraceKey = "stackTrace"

		// 创建encoder,同时输出到文件 + 控制台
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		// 动态设置日志级别
		level := zap.NewAtomicLevel()
		switch strings.ToLower(conf.Log.Level) {
		case "debug":
			level.SetLevel(zap.DebugLevel)
		case "info":
			level.SetLevel(zap.InfoLevel)
		case "warn":
			level.SetLevel(zap.WarnLevel)
		case "error":
			level.SetLevel(zap.ErrorLevel)
		case "dPanic":
			level.SetLevel(zap.DPanicLevel)
		case "panic":
			level.SetLevel(zap.PanicLevel)
		case "fatal":
			level.SetLevel(zap.FatalLevel)
		default:
			level.SetLevel(zap.InfoLevel)
		}

		// 创建核心
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
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

func (l *Logger) WithDebugger(c context.Context) *zap.Logger {
	var ms float64
	if start, ok := c.Value(ctxkey.StartTimeKey).(time.Time); ok {
		ms = float64(time.Since(start).Milliseconds())
	}
	if v, ok := c.Value(ctxkey.MsKey).(float64); ok {
		ms = v
	}

	traceId := c.Value(ctxkey.TraceIdKey).(string)

	return l.Logger.With(
		zap.String(ctxkey.TraceIdKey, traceId),
		zap.String(ctxkey.IpKey, c.Value(ctxkey.IpKey).(string)),
		zap.String(ctxkey.PathKey, c.Value(ctxkey.PathKey).(string)),
		zap.String(ctxkey.MethodKey, c.Value(ctxkey.MethodKey).(string)),
		zap.Any(ctxkey.ParamsKey, c.Value(ctxkey.ParamsKey)),
		zap.Float64(ctxkey.MsKey, ms),
		zap.Any(ctxkey.DebuggerKey, trace.Store.Get(traceId)),
	)
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

package facade

import (
	"context"
	"gin/pkg/logger"
	"go.uber.org/zap"
)

// Log 日志门面-日志记录统一入口
var Log = &logFacade{}

type logFacade struct{}

// Logger 获取日志实例(不带上下文)
func (l *logFacade) Logger() *logger.Logger {
	if log := Get("log"); log != nil {
		return log.(*logger.Logger)
	}
	return logger.NewLogger()
}

func (l *logFacade) WithDebugger(ctx context.Context) *zap.Logger {
	return l.Logger().WithDebugger(ctx)
}

func (l *logFacade) Debug(msg string, fields ...zap.Field) {
	l.Logger().Debug(msg, fields...)
}

func (l *logFacade) Info(msg string, fields ...zap.Field) {
	l.Logger().Info(msg, fields...)
}

func (l *logFacade) Warn(msg string, fields ...zap.Field) {
	l.Logger().Warn(msg, fields...)
}

func (l *logFacade) Error(msg string, fields ...zap.Field) {
	l.Logger().Error(msg, fields...)
}

func (l *logFacade) Fatal(msg string, fields ...zap.Field) {
	l.Logger().Fatal(msg, fields...)
}

func (l *logFacade) Panic(msg string, fields ...zap.Field) {
	l.Logger().Panic(msg, fields...)
}

func (l *logFacade) DPanic(msg string, fields ...zap.Field) {
	l.Logger().DPanic(msg, fields...)
}

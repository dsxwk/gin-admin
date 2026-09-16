package tests

import (
	"bytes"
	"encoding/json"
	"gin/app/facade"
	"gin/pkg/serviceprovider/logger"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志级别切换测试
func TestLoggerLevel(t *testing.T) {
	log := facade.Log()
	original := log.GetLevel()
	defer log.SetLevel(original)

	levels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
	for _, level := range levels {
		log.SetLevel(level)
		if got := log.GetLevel(); got != level {
			t.Fatalf("expected level %s, got %s", level, got)
		}
	}

	log.SetLevel("invalid")
	if got := log.GetLevel(); got != "info" {
		t.Fatalf("expected fallback level info, got %s", got)
	}
}

// 结构化堆栈测试
func TestStructuredStackCore(t *testing.T) {
	var output bytes.Buffer
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.StacktraceKey = logger.StackTraceKey

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)
	log := zap.New(
		logger.StructuredStackCore{Core: core},
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	log.Error("test")

	var data map[string]any
	if err := json.Unmarshal(output.Bytes(), &data); err != nil {
		t.Fatalf("unmarshal log failed: %v", err)
	}

	stack, ok := data[logger.StackTraceKey].([]any)
	if !ok {
		t.Fatalf("expected stackTrace array, got %T", data[logger.StackTraceKey])
	}
	if len(stack) == 0 {
		t.Fatal("expected stackTrace lines")
	}
}

// 控制台堆栈换行测试
func TestConsoleStackMultiline(t *testing.T) {
	var output bytes.Buffer
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.StacktraceKey = logger.StackTraceKey

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)
	log := zap.New(
		core,
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	log.Error("test")

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiline stack, got %q", output.String())
	}
}

package tests

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/debugger"
	"testing"
	"time"
)

// 调试器重复启动测试
func TestDebuggerStartIdempotent(t *testing.T) {
	bus := facade.Event().Bus()
	instance := facade.Debugger()
	defer func() {
		if !instance.IsRunning() {
			instance.Start()
		}
	}()

	instance.Start()
	if count := bus.Count(debugger.TopicSQL); count != 1 {
		t.Fatalf("expected one sql subscription, got %d", count)
	}

	instance.Start()
	if count := bus.Count(debugger.TopicSQL); count != 1 {
		t.Fatalf("expected idempotent start to keep one sql subscription, got %d", count)
	}

	instance.Stop()
	if instance.IsRunning() {
		t.Fatal("expected debugger stopped")
	}
	if count := bus.Count(debugger.TopicSQL); count != 0 {
		t.Fatalf("expected sql subscription removed, got %d", count)
	}

	instance.Start()
	if !instance.IsRunning() {
		t.Fatal("expected debugger restarted")
	}
}

// 追踪数据过期清理测试
func TestTraceCleanExpired(t *testing.T) {
	store := debugger.NewTrace()
	store.GetOrCreate("clean-expired")

	if cleaned := store.CleanExpired(time.Hour); cleaned != 0 {
		t.Fatalf("expected fresh trace not cleaned, got %d", cleaned)
	}

	time.Sleep(2 * time.Millisecond)
	if cleaned := store.CleanExpired(time.Nanosecond); cleaned != 1 {
		t.Fatalf("expected expired trace cleaned, got %d", cleaned)
	}
}

// 追踪数据写入测试
func TestDebuggerAddTraceEvent(t *testing.T) {
	store := debugger.NewTrace()
	trace := store.GetOrCreate("trace-event")
	trace.Append(debugger.SQLEvent{
		TraceID: "trace-event",
		SQL:     "select 1",
	})
	trace.Append(debugger.CacheEvent{
		TraceID: "trace-event",
		Name:    "user",
	})

	data, ok := store.Get("trace-event")
	if !ok {
		t.Fatal("expected trace data")
	}
	if got := len(data.SQL); got != 1 {
		t.Fatalf("expected one sql item, got %d", got)
	}
	if got := len(data.Cache); got != 1 {
		t.Fatalf("expected one cache item, got %d", got)
	}
}

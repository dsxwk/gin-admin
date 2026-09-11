package tests

import (
	"gin/app/facade"
	"testing"
)

// 事件总线发布和取消订阅测试
func TestBusPublishAndUnsubscribe(t *testing.T) {
	bus := facade.Event().Bus()
	bus.Clear()
	defer bus.Clear()

	topic := "test:bus"
	got := 0
	id := bus.Subscribe(topic, func(value int) {
		got = value
	})

	bus.Publish(topic, 1)
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	if !bus.Unsubscribe(topic, id) {
		t.Fatal("expected unsubscribe success")
	}

	bus.Publish(topic, 2)
	if got != 1 {
		t.Fatalf("expected 1 after unsubscribe, got %d", got)
	}
}

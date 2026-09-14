package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"testing"
	"time"
)

// 事件总线发布和取消订阅测试
func TestBusPublishAndUnsubscribe(t *testing.T) {
	bus := facade.Event().Bus()

	topic := "test:bus:unsubscribe"
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

// 事件总线异步发布测试
func TestBusAsyncPublish(t *testing.T) {
	bus := facade.Event().Bus()
	done := make(chan struct{})

	topic := "test:bus:async"
	id := bus.SubscribeAsync(topic, func(value int) {
		time.Sleep(20 * time.Millisecond)
		if value == 1 {
			close(done)
		}
	})
	defer bus.Unsubscribe(topic, id)

	bus.Publish(topic, 1)

	select {
	case <-done:
		t.Fatal("expected asynchronous publish")
	default:
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected async listener to complete")
	}
}

// 事件总线上下文取消测试
func TestBusPublishWithCanceledContext(t *testing.T) {
	bus := facade.Event().Bus()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	topic := "test:bus:cancel"
	called := false
	id := bus.Subscribe(topic, func(value int) {
		called = true
	})
	defer bus.Unsubscribe(topic, id)

	bus.PublishWithContext(ctx, topic, 1)
	if called {
		t.Fatal("expected canceled context to stop publish")
	}
}

type eventListTestEvent struct{}

func (e eventListTestEvent) Name() string {
	return "test:event-list"
}

func (e eventListTestEvent) Description() string {
	return "事件列表副本测试"
}

type eventListTestListener struct{}

func (l eventListTestListener) Handle(event eventListTestEvent) {}

// 初始化监听器注册测试
func TestInitListenerRegistered(t *testing.T) {
	list := facade.Event().List()
	for _, item := range list {
		if item.Name == "user.login" && len(item.Listeners) >= 2 {
			return
		}
	}

	t.Fatal("expected user.login init listeners")
}

// 事件列表监听器切片副本测试
func TestEventListReturnsListenerCopy(t *testing.T) {
	facade.Event().Register(eventListTestListener{}, eventListTestEvent{})

	list := facade.Event().List()
	targetIndex := -1
	for i := range list {
		if list[i].Name == "test:event-list" {
			targetIndex = i
			break
		}
	}
	if targetIndex < 0 || len(list[targetIndex].Listeners) == 0 {
		t.Fatal("expected registered event listener")
	}

	target := &list[targetIndex]
	target.Listeners[0] = "changed"

	list = facade.Event().List()
	for i := range list {
		if list[i].Name == "test:event-list" && list[i].Listeners[0] == "changed" {
			t.Fatal("expected event list listener slice copy")
		}
	}
}

// 业务事件收集测试
func TestBusinessEventCollected(t *testing.T) {
	traceID := "test-business-event"
	debugger.Store.Delete(traceID)
	defer debugger.Store.Delete(traceID)

	ctx := context.WithValue(context.Background(), ctxkey.TraceIdKey, traceID)
	facade.Event().Publish(ctx, eventListTestEvent{})

	trace, ok := debugger.Store.Get(traceID)
	if !ok {
		t.Fatal("expected business event trace")
	}
	if len(trace.Listener) != 1 {
		t.Fatalf("expected one listener event, got %d", len(trace.Listener))
	}
	if trace.Listener[0].Name != "test:event-list" {
		t.Fatalf("expected event name test:event-list, got %s", trace.Listener[0].Name)
	}
}

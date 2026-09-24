package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"testing"
	"time"
)

// 事件总线发布和取消订阅测试
func TestBusPublishAndUnsubscribe(t *testing.T) {
	bus := facade.Event().Bus()

	topic := "test:bus:unsubscribe"
	got := 0
	subscription := bus.Subscribe(topic, func(_ context.Context, value int) {
		got = value
	})

	bus.Publish(context.Background(), topic, 1)
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	if !subscription.Unsubscribe() {
		t.Fatal("expected unsubscribe success")
	}

	bus.Publish(context.Background(), topic, 2)
	if got != 1 {
		t.Fatalf("expected 1 after unsubscribe, got %d", got)
	}
}

// 事件总线异步发布测试
func TestBusAsyncPublish(t *testing.T) {
	bus := facade.Event().Bus()
	done := make(chan struct{})

	topic := "test:bus:async"
	subscription := bus.SubscribeAsync(topic, func(_ context.Context, value int) {
		time.Sleep(20 * time.Millisecond)
		if value == 1 {
			close(done)
		}
	})
	defer subscription.Unsubscribe()

	bus.Publish(context.Background(), topic, 1)

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

// 事件处理panic回调测试
func TestBusPanicHandler(t *testing.T) {
	called := make(chan struct{})
	bus := eventbus.NewBus(
		eventbus.WithWorkers(1),
		eventbus.WithQueueSize(1),
		eventbus.WithPanicHandler(func(topic string, event any, recovered any) {
			close(called)
		}),
	)
	t.Cleanup(func() { _ = bus.Close(context.Background()) })

	subscription := bus.SubscribeAsync("test:bus:panic", func(_ context.Context, _ int) {
		panic("test panic")
	})
	defer subscription.Unsubscribe()

	bus.Publish(context.Background(), "test:bus:panic", 1)

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("expected panic handler called")
	}
}

// 事件总线上下文取消测试
func TestBusPublishWithCanceledContext(t *testing.T) {
	bus := facade.Event().Bus()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	topic := "test:bus:cancel"
	called := false
	subscription := bus.Subscribe(topic, func(_ context.Context, value int) {
		called = true
	})
	defer subscription.Unsubscribe()

	bus.Publish(ctx, topic, 1)
	if called {
		t.Fatal("expected canceled context to stop publish")
	}
}

// 异步事件不受发布方context取消影响
func TestBusAsyncSurvivesContextCancel(t *testing.T) {
	bus := eventbus.NewBus(
		eventbus.WithWorkers(1),
		eventbus.WithQueueSize(8),
	)
	t.Cleanup(func() { _ = bus.Close(context.Background()) })

	traceID := "test-bus-async-cancel"
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxkey.TraceIDKey, traceID))

	done := make(chan struct{})
	topic := "test:bus:async-cancel"
	subscription := bus.SubscribeAsync(topic, func(handlerCtx context.Context, _ int) {
		if handlerCtx.Err() != nil {
			t.Error("expected handler context not canceled")
		}
		if value, _ := handlerCtx.Value(ctxkey.TraceIDKey).(string); value != traceID {
			t.Errorf("expected traceID %s, got %s", traceID, value)
		}
		close(done)
	})
	defer subscription.Unsubscribe()

	bus.Publish(ctx, topic, 1)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected async handler run after context canceled")
	}
}

// 队列满丢弃事件且不阻塞发布方
func TestBusQueueFullDrops(t *testing.T) {
	bus := eventbus.NewBus(
		eventbus.WithWorkers(1),
		eventbus.WithQueueSize(1),
		eventbus.WithEnqueueWait(10*time.Millisecond),
	)
	t.Cleanup(func() { _ = bus.Close(context.Background()) })

	started := make(chan struct{}, 8)
	release := make(chan struct{})
	topic := "test:bus:queue-full"
	subscription := bus.SubscribeAsync(topic, func(_ context.Context, _ int) {
		started <- struct{}{}
		<-release
	})
	defer subscription.Unsubscribe()

	bus.Publish(context.Background(), topic, 1)

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("expected handler started")
	}

	// worker 被占用,队列已被第二个事件填满
	bus.Publish(context.Background(), topic, 2)

	begin := time.Now()
	bus.Publish(context.Background(), topic, 3)
	if elapsed := time.Since(begin); elapsed > time.Second {
		t.Fatalf("expected publish not blocked, elapsed %s", elapsed)
	}

	if dropped := bus.Stats().Dropped; dropped != 1 {
		t.Fatalf("expected 1 dropped event, got %d", dropped)
	}

	close(release)
}

// 事件类型不匹配计数
func TestBusTypeMismatchCounted(t *testing.T) {
	bus := eventbus.NewBus(
		eventbus.WithWorkers(1),
		eventbus.WithQueueSize(4),
	)
	t.Cleanup(func() { _ = bus.Close(context.Background()) })

	topic := "test:bus:mismatch"
	subscription := bus.SubscribeAsync(topic, func(_ context.Context, _ int) {})
	defer subscription.Unsubscribe()

	bus.Publish(context.Background(), topic, "not-int")

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if bus.Stats().Mismatched == 1 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}

	t.Fatalf("expected 1 mismatched event, got %d", bus.Stats().Mismatched)
}

type eventListTestEvent struct{}

func (e eventListTestEvent) Name() string {
	return "test:event-list"
}

func (e eventListTestEvent) Description() string {
	return "事件列表副本测试"
}

type eventListTestListener struct{}

func (l eventListTestListener) Handle(_ eventListTestEvent) {}

type duplicateEventTestEvent struct{}

func (e duplicateEventTestEvent) Name() string {
	return "test:event-duplicate"
}

func (e duplicateEventTestEvent) Description() string {
	return "重复监听器测试"
}

type duplicateEventTestListener struct{}

func (l duplicateEventTestListener) Handle(_ duplicateEventTestEvent) {}

type emptyEventTestEvent struct{}

func (e emptyEventTestEvent) Name() string {
	return ""
}

func (e emptyEventTestEvent) Description() string {
	return "空事件名称测试"
}

type emptyEventTestListener struct{}

func (l emptyEventTestListener) Handle(_ emptyEventTestEvent) {}

// 默认监听器注册测试
func TestDefaultListenerRegistered(t *testing.T) {
	list := facade.Event().List()
	for _, item := range list {
		if item.Name == "user.login" && len(item.Listeners) >= 2 {
			return
		}
	}

	t.Fatal("expected user.login listeners")
}

// 重复监听器注册测试
func TestRegistryRegisterDuplicate(t *testing.T) {
	bus := eventbus.NewBus()
	t.Cleanup(func() { _ = bus.Close(context.Background()) })
	registry := eventbus.NewRegistry(bus)
	listener := duplicateEventTestListener{}
	e := duplicateEventTestEvent{}

	if !registry.Register(listener, e) {
		t.Fatal("expected first listener registered")
	}
	if registry.Register(listener, e) {
		t.Fatal("expected duplicate listener skipped")
	}

	list := registry.EventList()
	if len(list) != 1 || len(list[0].Listeners) != 1 {
		t.Fatalf("expected one listener, got %+v", list)
	}
}

// 空事件名称注册测试
func TestRegistryRegisterEmptyName(t *testing.T) {
	bus := eventbus.NewBus()
	t.Cleanup(func() { _ = bus.Close(context.Background()) })
	registry := eventbus.NewRegistry(bus)

	if registry.Register(emptyEventTestListener{}, emptyEventTestEvent{}) {
		t.Fatal("expected empty event name rejected")
	}
	if len(registry.EventList()) != 0 {
		t.Fatal("expected empty event not registered")
	}
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

	ctx := context.WithValue(context.Background(), ctxkey.TraceIDKey, traceID)
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

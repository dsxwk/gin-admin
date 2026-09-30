package ws

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

const maxInitialShardCapacity = 4096

// shard WebSocket连接分片
type shard struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

// Manager WebSocket连接管理器
type Manager struct {
	options        Options
	handlers       []Handler
	upgrader       websocket.Upgrader
	shards         []*shard
	nextClientID   atomic.Uint64
	connectionSize atomic.Int64
	lifecycleMu    sync.RWMutex
	closed         bool
	closeOnce      sync.Once
	wg             sync.WaitGroup
	shutdownWG     sync.WaitGroup
}

// NewManager 创建WebSocket管理器
func NewManager(options Options, handlers ...Handler) *Manager {
	options = options.normalize()
	manager := &Manager{
		options:  options,
		handlers: make([]Handler, 0, len(handlers)),
		shards:   make([]*shard, options.Shards),
		upgrader: websocket.Upgrader{
			ReadBufferSize:    options.ReadBufferSize,
			WriteBufferSize:   options.WriteBufferSize,
			EnableCompression: options.Compression,
		},
	}

	for _, handler := range handlers {
		if handler != nil {
			manager.handlers = append(manager.handlers, handler)
		}
	}

	for index := range manager.shards {
		capacity := min(options.MaxConnections/int64(options.Shards)+1, maxInitialShardCapacity)
		manager.shards[index] = &shard{
			clients: make(map[string]*Client, int(capacity)),
		}
	}
	manager.upgrader.CheckOrigin = manager.checkOrigin

	return manager
}

// ServeHTTP 处理WebSocket升级请求
func (m *Manager) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if m == nil {
		http.Error(writer, "websocket服务未初始化", http.StatusServiceUnavailable)
		return
	}

	m.lifecycleMu.RLock()
	if m.closed {
		m.lifecycleMu.RUnlock()
		http.Error(writer, ErrClosed.Error(), http.StatusServiceUnavailable)
		return
	}
	if !m.reserveConnection() {
		m.lifecycleMu.RUnlock()
		http.Error(writer, ErrTooManyConnections.Error(), http.StatusServiceUnavailable)
		return
	}

	conn, err := m.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		m.connectionSize.Add(-1)
		m.lifecycleMu.RUnlock()
		return
	}

	clientID := m.nextClientID.Add(1)
	client := newClient(m, m.shardFor(clientID), conn, request.Context(), clientID)
	client.BindUser(UserID(request.Context()))
	client.shard.add(client)
	m.wg.Add(2)
	m.lifecycleMu.RUnlock()

	go client.writePump()
	go client.readPump()
}

// reserveConnection 预占连接数
func (m *Manager) reserveConnection() bool {
	for {
		current := m.connectionSize.Load()
		if current >= m.options.MaxConnections {
			return false
		}
		if m.connectionSize.CompareAndSwap(current, current+1) {
			return true
		}
	}
}

// shardFor 获取客户端分片
func (m *Manager) shardFor(clientID uint64) *shard {
	return m.shards[clientID%uint64(len(m.shards))]
}

// unregister 移除客户端连接
func (m *Manager) unregister(client *Client) {
	if client == nil || client.shard == nil {
		return
	}

	client.shard.remove(client.id)
	m.connectionSize.Add(-1)
	close(client.done)
	_ = client.conn.Close()
}

// dispatch 分发消息并在处理器panic时保护服务
func (m *Manager) dispatch(ctx context.Context, client *Client, message Message) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("websocket处理器panic: %v", recovered)
		}
	}()

	hasMatchedHandler := false
	for _, handler := range m.handlers {
		matcher, ok := handler.(HandlerMatcher)
		if ok && matcher.Match(message) {
			hasMatchedHandler = true
			break
		}
	}

	if hasMatchedHandler {
		for _, handler := range m.handlers {
			matcher, ok := handler.(HandlerMatcher)
			if !ok || !matcher.Match(message) {
				continue
			}
			if err = handler.Handle(ctx, client, message); err != nil {
				return err
			}
		}
		return nil
	}

	for _, handler := range m.handlers {
		if _, ok := handler.(HandlerMatcher); ok {
			continue
		}
		if err = handler.Handle(ctx, client, message); err != nil {
			return err
		}
	}

	return nil
}

// checkOrigin 校验WebSocket来源
func (m *Manager) checkOrigin(request *http.Request) bool {
	origin := strings.TrimSpace(request.Header.Get("Origin"))
	if origin == "" {
		return true
	}

	for _, allowed := range m.options.AllowedOrigins {
		allowed = strings.TrimRight(strings.TrimSpace(allowed), "/")
		if allowed == "*" || strings.EqualFold(allowed, strings.TrimRight(origin, "/")) {
			return true
		}
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsed.Host, request.Host)
}

// Client 获取指定客户端
func (m *Manager) Client(clientID string) *Client {
	if m == nil || clientID == "" {
		return nil
	}

	return m.shards[m.shardIndex(clientID)].get(clientID)
}

// Count 获取当前连接数
func (m *Manager) Count() int64 {
	if m == nil {
		return 0
	}
	return m.connectionSize.Load()
}

// Send 发送消息到指定客户端
func (m *Manager) Send(clientID string, message Message) error {
	client := m.Client(clientID)
	if client == nil {
		return ErrClosed
	}

	return client.Send(message)
}

// SendToUser 发送消息到用户全部连接
func (m *Manager) SendToUser(userID int64, message Message) int {
	if m == nil || userID <= 0 {
		return 0
	}

	message = cloneMessage(message)
	counts := make([]int, len(m.shards))
	var waitGroup sync.WaitGroup

	for index, currentShard := range m.shards {
		waitGroup.Go(func() {
			counts[index] = currentShard.sendToUser(userID, message)
		})
	}
	waitGroup.Wait()

	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

// SendToUsers 发送消息到多个用户全部连接
func (m *Manager) SendToUsers(userIDs []int64, message Message) int {
	if m == nil || len(userIDs) == 0 {
		return 0
	}

	targets := make(map[int64]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if userID > 0 {
			targets[userID] = struct{}{}
		}
	}
	if len(targets) == 0 {
		return 0
	}

	message = cloneMessage(message)
	counts := make([]int, len(m.shards))
	var waitGroup sync.WaitGroup

	for index, currentShard := range m.shards {
		waitGroup.Go(func() {
			counts[index] = currentShard.sendToUsers(targets, message)
		})
	}
	waitGroup.Wait()

	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

// Broadcast 广播消息
func (m *Manager) Broadcast(message Message) int {
	if m == nil {
		return 0
	}

	message = cloneMessage(message)
	counts := make([]int, len(m.shards))
	var waitGroup sync.WaitGroup

	for index, currentShard := range m.shards {
		waitGroup.Go(func() {
			counts[index] = currentShard.broadcast(message)
		})
	}
	waitGroup.Wait()

	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

// Stats WebSocket运行统计
type Stats struct {
	Connections    int64 // 当前连接数
	MaxConnections int64 // 最大连接数
	Shards         int   // 分片数量
}

// Stats 获取WebSocket运行统计
func (m *Manager) Stats() Stats {
	if m == nil {
		return Stats{}
	}

	return Stats{
		Connections:    m.connectionSize.Load(),
		MaxConnections: m.options.MaxConnections,
		Shards:         len(m.shards),
	}
}

// Close 关闭全部WebSocket连接
func (m *Manager) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	m.closeOnce.Do(func() {
		m.lifecycleMu.Lock()
		m.closed = true
		m.lifecycleMu.Unlock()

		for _, currentShard := range m.shards {
			m.shutdownWG.Go(func() {
				for _, client := range currentShard.list() {
					client.shutdown()
				}
			})
		}
	})

	done := make(chan struct{})
	go func() {
		m.shutdownWG.Wait()
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// shardIndex 计算分片下标
func (m *Manager) shardIndex(clientID string) uint64 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(clientID))
	return uint64(hash.Sum32()) % uint64(len(m.shards))
}

// add 添加客户端
func (s *shard) add(client *Client) {
	s.mu.Lock()
	s.clients[client.id] = client
	s.mu.Unlock()
}

// remove 移除客户端
func (s *shard) remove(clientID string) {
	s.mu.Lock()
	delete(s.clients, clientID)
	s.mu.Unlock()
}

// get 获取客户端
func (s *shard) get(clientID string) *Client {
	s.mu.RLock()
	client := s.clients[clientID]
	s.mu.RUnlock()
	return client
}

// list 获取分片客户端列表
func (s *shard) list() []*Client {
	s.mu.RLock()
	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	s.mu.RUnlock()
	return clients
}

// broadcast 广播分片消息
func (s *shard) broadcast(message Message) int {
	clients := s.list()
	count := 0
	for _, client := range clients {
		if client.enqueue(message) == nil {
			count++
		}
	}
	return count
}

// sendToUser 发送消息到分片内指定用户
func (s *shard) sendToUser(userID int64, message Message) int {
	clients := s.list()
	count := 0
	for _, client := range clients {
		if client.UserID() != userID {
			continue
		}
		if client.enqueue(message) == nil {
			count++
		}
	}
	return count
}

// sendToUsers 发送消息到分片内多个用户
func (s *shard) sendToUsers(userIDs map[int64]struct{}, message Message) int {
	clients := s.list()
	count := 0
	for _, client := range clients {
		if _, ok := userIDs[client.UserID()]; !ok {
			continue
		}
		if client.enqueue(message) == nil {
			count++
		}
	}
	return count
}

package ws

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端连接
type Client struct {
	id        string
	manager   *Manager
	shard     *shard
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	userID    atomic.Int64
	sendQueue chan Message
	done      chan struct{}
	closeOnce sync.Once
}

// newClient 创建WebSocket客户端
func newClient(manager *Manager, shard *shard, conn *websocket.Conn, requestContext context.Context, clientID uint64) *Client {
	if requestContext == nil {
		requestContext = context.Background()
	}

	ctx, cancel := context.WithCancel(context.WithoutCancel(requestContext))

	return &Client{
		id:        strconv.FormatUint(clientID, 10),
		manager:   manager,
		shard:     shard,
		conn:      conn,
		ctx:       ctx,
		cancel:    cancel,
		sendQueue: make(chan Message, manager.options.WriteQueueSize),
		done:      make(chan struct{}),
	}
}

// ID 获取客户端ID
func (c *Client) ID() string {
	if c == nil {
		return ""
	}
	return c.id
}

// UserID 获取用户ID
func (c *Client) UserID() int64 {
	if c == nil {
		return 0
	}
	return c.userID.Load()
}

// BindUser 绑定用户
func (c *Client) BindUser(userID int64) {
	if c == nil || userID < 0 {
		return
	}
	if current := c.userID.Load(); current > 0 && current != userID {
		return
	}
	c.userID.Store(userID)
}

// Send 发送WebSocket消息
func (c *Client) Send(message Message) error {
	if c == nil {
		return ErrClosed
	}

	return c.enqueue(cloneMessage(message))
}

// enqueue 加入发送队列
func (c *Client) enqueue(message Message) error {
	select {
	case <-c.done:
		return ErrClosed
	default:
	}

	select {
	case c.sendQueue <- message:
		return nil
	case <-c.done:
		return ErrClosed
	default:
		c.Close()
		return ErrSlowClient
	}
}

// Close 关闭WebSocket连接
func (c *Client) Close() {
	if c == nil {
		return
	}

	c.closeOnce.Do(func() {
		c.cancel()
		c.manager.unregister(c)
	})
}

// shutdown 发送关闭帧并关闭连接
func (c *Client) shutdown() {
	if c == nil {
		return
	}

	_ = c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutdown"),
		time.Now().Add(c.manager.options.WriteWait),
	)
	c.Close()
}

// readPump 读取WebSocket消息
func (c *Client) readPump() {
	defer c.manager.wg.Done()
	defer c.Close()

	c.conn.SetReadLimit(c.manager.options.ReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.manager.options.PongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(c.manager.options.PongWait))
	})

	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		data = append([]byte(nil), data...)
		err = c.manager.dispatch(c.ctx, c, Message{
			Type: MessageType(messageType),
			Data: data,
		})
		if err != nil {
			_ = c.Send(Message{
				Type: MessageText,
				Data: []byte(err.Error()),
			})
		}
	}
}

// writePump 写入WebSocket消息
func (c *Client) writePump() {
	ticker := time.NewTicker(c.manager.options.PingInterval)
	defer func() {
		ticker.Stop()
		c.manager.wg.Done()
		c.Close()
	}()

	for {
		select {
		case message := <-c.sendQueue:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.manager.options.WriteWait))
			if err := c.conn.WriteMessage(int(message.Type), message.Data); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(c.manager.options.WriteWait),
			); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}

// cloneMessage 复制消息内容
func cloneMessage(message Message) Message {
	message.Data = append([]byte(nil), message.Data...)
	return message
}

package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client grpc客户端
type Client struct {
	Conn *grpc.ClientConn
}

// NewClient 创建grpc客户端
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(unaryClientInterceptor),
	)
	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}

// Service 获取指定服务客户端
func (c *Client) Service[T any](factory func(grpc.ClientConnInterface) T) T {
	return factory(c.Conn)
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

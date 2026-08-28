package grpcclient

import (
	"errors"
	"fmt"
	"gin/common/flag"
	"net"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server grpc服务端
type Server struct {
	server   *grpclib.Server
	listener net.Listener
	addr     string
}

// NewServer 创建grpc服务端
func NewServer(host string, port int, registrars ...ServiceRegistrar) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	s := grpclib.NewServer(
		grpclib.ChainUnaryInterceptor(unaryServerInterceptor),
	)
	reflection.Register(s)
	for _, register := range Registrars() {
		register(s)
	}
	for _, register := range registrars {
		register(s)
	}

	return &Server{
		server:   s,
		listener: listener,
		addr:     listener.Addr().String(),
	}, nil
}

// Addr 获取监听地址
func (s *Server) Addr() string {
	return s.addr
}

// Start 启动服务
func (s *Server) Start() error {
	go func() {
		if err := s.server.Serve(s.listener); err != nil && !errors.Is(err, grpclib.ErrServerStopped) {
			flag.Errorf("grpc服务停止: %v", err)
		}
	}()
	return nil
}

// Stop 停止服务
func (s *Server) Stop() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

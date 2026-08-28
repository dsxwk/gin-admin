package grpcclient

import (
	"sync"

	grpclib "google.golang.org/grpc"
)

// ServiceRegistrar grpc服务注册回调
type ServiceRegistrar func(s *grpclib.Server)

var (
	serviceRegistrars  []ServiceRegistrar
	serviceRegistrarMu sync.RWMutex
)

// Register 注册grpc服务
func Register(registrar ServiceRegistrar) {
	serviceRegistrarMu.Lock()
	defer serviceRegistrarMu.Unlock()
	serviceRegistrars = append(serviceRegistrars, registrar)
}

// Registrars 获取已注册的grpc服务
func Registrars() []ServiceRegistrar {
	serviceRegistrarMu.RLock()
	defer serviceRegistrarMu.RUnlock()
	return append([]ServiceRegistrar(nil), serviceRegistrars...)
}

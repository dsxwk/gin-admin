package route

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	registryMu      sync.RWMutex
	registryRouters []Router
)

// Register 注册路由模块
func Register(r Router) {
	if r == nil {
		return
	}

	registryMu.Lock()
	registryRouters = append(registryRouters, r)
	registryMu.Unlock()
}

// Mount 挂载已注册路由
func Mount(public, auth *gin.RouterGroup) {
	for _, r := range snapshot() {
		if r.IsAuth() {
			r.RegisterRoutes(auth)
			continue
		}

		r.RegisterRoutes(public)
	}
}

// snapshot 获取路由快照
func snapshot() []Router {
	registryMu.RLock()
	defer registryMu.RUnlock()

	routers := make([]Router, len(registryRouters))
	copy(routers, registryRouters)

	return routers
}

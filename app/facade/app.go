package facade

import "gin/pkg/container"

// App 获取服务容器
func App() *container.Container {
	return container.Default()
}

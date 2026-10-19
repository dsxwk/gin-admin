package facade

import (
	serviceagent "gin/pkg/serviceprovider/agent"
)

// Agent 获取AI Agent管理器
func Agent() *serviceagent.Manager {
	return App().Agent()
}

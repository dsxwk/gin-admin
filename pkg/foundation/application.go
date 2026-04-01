package foundation

import (
	"context"
	"fmt"
	"sync"
)

// Application 应用实现
type Application struct {
	mu           sync.RWMutex
	providers    []ServiceProvider  // 已注册的服务提供者列表
	runners      []Runner           // 需要后台运行的任务
	runnerCancel context.CancelFunc // 后台任务取消函数
	runnerCtx    context.Context    // 后台任务上下文
	initialized  bool               // 应用是否已启动
}

var (
	defaultApp *Application
	appOnce    sync.Once
)

// NewApplication 创建应用实例
func NewApplication() *Application {
	return &Application{
		providers: make([]ServiceProvider, 0),
		runners:   make([]Runner, 0),
	}
}

// GetApp 获取应用单例(自动初始化)
func GetApp() *Application {
	appOnce.Do(func() {
		defaultApp = NewApplication()
	})
	return defaultApp
}

// Register 注册服务提供者
func (app *Application) Register(providers ...ServiceProvider) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.providers = append(app.providers, providers...)
}

// Boot 启动应用,自动加载所有通过init注册的服务提供者
func (app *Application) Boot() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.initialized {
		return nil
	}

	// 获取所有自动注册的服务提供者
	providers := GetProviders()
	if len(providers) == 0 {
		return nil
	}

	// 按依赖关系排序
	// 确保依赖的服务先启动,例如:
	//   DatabaseProvider依赖ConfigProvider→ConfigProvider先启动
	//   CacheProvider依赖ConfigProvider→ConfigProvider先启动
	//   RateLimitProvider依赖ConfigProvider和CacheProvider→两者都先启动
	sortedProviders, err := app.sortProvidersByDependency(providers)
	if err != nil {
		return fmt.Errorf("failed to sort providers: %w", err)
	}

	// 调用所有服务的Register方法
	for _, provider := range sortedProviders {
		provider.Register(app)
	}

	// 按顺序调用所有服务的Boot方法
	for _, provider := range sortedProviders {
		provider.Boot(app)

		// 收集后台运行任务
		if withRunners, ok := provider.(ServiceProviderWithRunners); ok {
			app.runners = append(app.runners, withRunners.Runners()...)
		}
	}

	// 启动所有后台任务
	app.startRunners()
	// 标记为已初始化,防止重复启动
	app.initialized = true

	return nil
}

// Stop 停止应用
func (app *Application) Stop() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.runnerCancel != nil {
		app.runnerCancel()
	}

	var lastErr error
	for _, runner := range app.runners {
		if err := runner.Stop(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// startRunners 启动所有后台运行任务
func (app *Application) startRunners() {
	if len(app.runners) == 0 {
		return
	}

	app.runnerCtx, app.runnerCancel = context.WithCancel(context.Background())

	for _, runner := range app.runners {
		go func(r Runner) {
			_ = r.Run(app.runnerCtx)
		}(runner)
	}
}

// sortProvidersByDependency 根据依赖关系排序服务提供者(拓扑排序)
// 确定服务提供者的启动顺序,确保依赖的服务先启动,被依赖的服务后启动
// 参数:
//   - providers: 待排序的服务提供者列表
//
// 返回:
//   - []ServiceProvider: 排序后的服务提供者列表
//   - error: 如果检测到循环依赖则返回错误
//
// 示例:
//
//	假设有3个服务: A(依赖B), B(依赖C), C(无依赖)
//	排序结果: [C, B, A]
func (app *Application) sortProvidersByDependency(providers []ServiceProvider) ([]ServiceProvider, error) {
	if len(providers) == 0 {
		return providers, nil
	}

	// 构建名称到索引的映射
	// {"config":0, "log":1, "database":2}
	nameToIndex := make(map[string]int)
	for i, p := range providers {
		nameToIndex[p.Name()] = i
	}

	// 构建依赖关系图
	// graph[i] = []int 表示服 i被哪些服务依赖
	// inDegree[i] 表示服务i依赖多少个其他服务
	//
	// 依赖关系示例:
	//   服务A依赖服务B和C
	//   则graph[B]和graph[C]包含A
	//   inDegree[A] = 2
	graph := make([][]int, len(providers))
	inDegree := make([]int, len(providers))

	// 遍历所有服务提供者,收集依赖关系
	for i, provider := range providers {
		// 检查服务提供者是否实现了依赖接口
		if withDeps, ok := provider.(ServiceProviderWithDependencies); ok {
			// 遍历该服务的所有依赖
			for _, depName := range withDeps.Dependencies() {
				// 查找依赖服务的索引
				if idx, exists := nameToIndex[depName]; exists {
					// idx->i(idx依赖的服务依赖于i)
					graph[idx] = append(graph[idx], i)
					// 依赖的服务数量
					inDegree[i]++
				}
			}
		}
	}

	// 拓扑排序
	// 使用Kahn算法进行拓扑排序
	// 找到所有入度为0的节点(没有依赖的服务)
	// 移除节点,并减少它们指向的节点的入度
	// 重复直到所有节点都被处理
	var result []ServiceProvider
	queue := make([]int, 0)

	// 找到所有入度为0的服务(没有依赖,可以最先启动)
	for i, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, i)
		}
	}

	// 循环处理队列
	for len(queue) > 0 {
		// 取出首个元素
		idx := queue[0]
		queue = queue[1:]
		result = append(result, providers[idx])

		// 遍历该服务被哪些服务依赖
		for _, neighbor := range graph[idx] {
			// 减少邻居的入度(因为依赖的服务已经被处理了)
			inDegree[neighbor]--
			// 如果邻居的入度变为0,说明它的所有依赖都已处理,可以加入队列
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 检测循环依赖
	// 如果排序后的数量不等于原始数量,说明存在循环依赖
	// 例如: A依赖B,B依赖A,则两者都无法被处理
	if len(result) != len(providers) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

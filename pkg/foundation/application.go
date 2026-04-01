package foundation

import (
	"context"
	"fmt"
	"sync"
)

// Application 应用实现
type Application struct {
	mu           sync.RWMutex
	providers    []ServiceProvider
	runners      []Runner
	runnerCancel context.CancelFunc
	runnerCtx    context.Context
	initialized  bool
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

// Boot 启动应用 - 自动加载所有通过 init 注册的服务提供者
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
func (app *Application) sortProvidersByDependency(providers []ServiceProvider) ([]ServiceProvider, error) {
	if len(providers) == 0 {
		return providers, nil
	}

	// 构建名称到索引的映射
	nameToIndex := make(map[string]int)
	for i, p := range providers {
		nameToIndex[p.Name()] = i
	}

	// 构建依赖图
	graph := make([][]int, len(providers))
	inDegree := make([]int, len(providers))

	for i, provider := range providers {
		if withDeps, ok := provider.(ServiceProviderWithDependencies); ok {
			for _, depName := range withDeps.Dependencies() {
				if idx, exists := nameToIndex[depName]; exists {
					graph[idx] = append(graph[idx], i)
					inDegree[i]++
				}
			}
		}
	}

	// 拓扑排序
	var result []ServiceProvider
	queue := make([]int, 0)

	for i, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]
		result = append(result, providers[idx])

		for _, neighbor := range graph[idx] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(result) != len(providers) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

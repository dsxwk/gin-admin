package serviceprovider

import (
	"context"
	"fmt"
	"gin/pkg/container"
	"sync"
	"time"
)

// Application 应用实现
type Application struct {
	mu           sync.RWMutex
	container    *container.Container
	providers    []ServiceProvider  // 服务提供者列表
	runners      []Runner           // 需要后台运行的任务
	runnerCancel context.CancelFunc // 后台任务取消函数
	registered   bool               // 服务是否已注册
	initialized  bool               // 应用是否已启动
}

// NewApp 创建应用实例
func NewApp(providers ...ServiceProvider) *Application {
	return &Application{
		container: container.Default(),
		providers: append([]ServiceProvider(nil), providers...),
		runners:   make([]Runner, 0),
	}
}

// RegisterProviders 仅执行服务注册阶段,不启动后台任务
func (app *Application) RegisterProviders() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	return app.registerProviders()
}

// Boot 启动应用
func (app *Application) Boot() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.initialized {
		return nil
	}

	if err := app.registerProviders(); err != nil {
		return err
	}

	for _, provider := range app.providers {
		provider.Boot(app.container)

		if withRunners, ok := provider.(Runners); ok {
			app.runners = append(app.runners, withRunners.Runners()...)
		}
	}

	app.startRunners()
	app.initialized = true

	return nil
}

// registerProviders 执行服务提供者注册阶段
func (app *Application) registerProviders() error {
	if app.registered {
		return nil
	}

	providers := app.providers
	if len(providers) == 0 {
		app.registered = true
		return nil
	}

	sortedProviders, err := app.sortProvidersByDependency(providers)
	if err != nil {
		return fmt.Errorf("failed to sort providers: %w", err)
	}

	for _, provider := range sortedProviders {
		provider.Register(app.container)
	}

	app.providers = sortedProviders
	app.registered = true
	return nil
}

// Stop 停止应用
func (app *Application) Stop() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.runnerCancel != nil {
		app.runnerCancel()
	}

	var errs []error
	for index := len(app.runners) - 1; index >= 0; index-- {
		runner := app.runners[index]
		if err := runner.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", runner.Name(), err))
		}
	}

	if registry := app.container.Event(); registry != nil {
		if bus := registry.Bus(); bus != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := bus.Close(ctx); err != nil {
				errs = append(errs, fmt.Errorf("eventbus: %w", err))
			}
			cancel()
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("stop errors: %v", errs)
	}
	return nil
}

// startRunners 启动所有后台运行任务
func (app *Application) startRunners() {
	if len(app.runners) == 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	app.runnerCancel = cancel

	for _, runner := range app.runners {
		go func(r Runner) {
			_ = r.Run(ctx)
		}(runner)
	}
}

// sortProvidersByDependency 根据依赖关系对服务提供者进行拓扑排序
// 返回顺序保证被依赖的Provider先执行
// 例如A依赖B,则结果一定是B在A之前
// 算法使用Kahn拓扑排序
// graph保存被依赖Provider指向依赖它的Provider的边
// inDegree保存每个Provider尚未处理的依赖数量
func (app *Application) sortProvidersByDependency(providers []ServiceProvider) ([]ServiceProvider, error) {
	if len(providers) == 0 {
		// 没有Provider时直接返回
		return providers, nil
	}

	// 建立Provider名称到切片下标的映射
	// 后续通过依赖名称快速定位Provider
	// 同时检测重复Provider避免依赖关系产生歧义
	nameToIndex := make(map[string]int)
	for i, p := range providers {
		name := p.Name()
		if _, exists := nameToIndex[name]; exists {
			return nil, fmt.Errorf("duplicate provider: %s", name)
		}
		nameToIndex[name] = i
	}

	// graph[idx]保存依赖idx的所有Provider下标
	// 例如A依赖B,则graph[B]包含A的下标
	// 这样处理完B后可以继续处理A
	graph := make([][]int, len(providers))

	// inDegree[i]表示Provider i还依赖多少个Provider
	// inDegree等于0说明Provider i没有未满足的依赖,可以立即执行
	inDegree := make([]int, len(providers))

	// 遍历全部Provider并建立依赖关系
	for i, provider := range providers {
		// 只有实现依赖接口的Provider才需要解析依赖
		if withDeps, ok := provider.(Dependencies); ok {
			for _, depName := range withDeps.Dependencies() {
				// 找不到依赖名称时忽略该依赖
				// 这样允许Provider声明可选依赖
				if idx, exists := nameToIndex[depName]; exists {
					// 增加一条idx指向i的边
					graph[idx] = append(graph[idx], i)

					// 记录i存在一个尚未处理的依赖
					inDegree[i]++
				}
			}
		}
	}

	// result 保存最终排序结果
	var result []ServiceProvider

	// queue 保存所有当前入度为0的Provider
	// 可以先用队列收集所有无依赖的Provider
	queue := make([]int, 0)

	// 将入度为0的Provider加入初始队列
	for i, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, i)
		}
	}

	// Kahn算法循环处理队列
	for len(queue) > 0 {
		// 取出队首Provider并加入结果
		idx := queue[0]
		queue = queue[1:]
		result = append(result, providers[idx])

		// 当前Provider已经处理完成
		// 因此依赖它的Provider可以减少一个未处理依赖
		for _, neighbor := range graph[idx] {
			inDegree[neighbor]--

			// 入度变为0说明依赖已经全部满足,可以加入队列
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 如果结果数量不等于Provider总数,说明存在循环依赖
	// 循环依赖中的Provider永远不会进入入度为0的队列
	if len(result) != len(providers) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

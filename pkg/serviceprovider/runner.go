package serviceprovider

import "context"

// shutdownRunner 通用关闭任务
type shutdownRunner struct {
	name string
	stop func() error
}

// NewShutdownRunner 创建通用关闭任务
func NewShutdownRunner(name string, stop func() error) Runner {
	return &shutdownRunner{
		name: name,
		stop: stop,
	}
}

// Run 等待应用停止
func (r *shutdownRunner) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	<-ctx.Done()
	return nil
}

// Stop 执行关闭回调
func (r *shutdownRunner) Stop() error {
	if r.stop == nil {
		return nil
	}
	return r.stop()
}

// Name 获取任务名称
func (r *shutdownRunner) Name() string {
	return r.name
}

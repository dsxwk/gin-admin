package facade

import (
	"context"
	"fmt"
	"gin/pkg/container"
	jsjob "gin/pkg/serviceprovider/job"
)

// Job 任务门面实例
func Job() *JobFacade {
	return &JobFacade{manager: container.Default().Job()}
}

// JobFacade 任务门面
type JobFacade struct {
	manager *jsjob.Manager
}

func (j *JobFacade) Dispatch(ctx context.Context, jobName string, payload any) error {
	if j == nil || j.manager == nil {
		return fmt.Errorf("job manager 未初始化")
	}
	return j.manager.Dispatch(ctx, jobName, payload)
}

// Jobs 获取所有Job
func (j *JobFacade) Jobs() []jsjob.Job {
	if j == nil || j.manager == nil {
		return nil
	}
	return j.manager.Jobs()
}

func (j *JobFacade) Count(ctx context.Context) (int64, error) {
	if j == nil || j.manager == nil {
		return 0, fmt.Errorf("job manager 未初始化")
	}
	return j.manager.Count(ctx)
}

func (j *JobFacade) Clear(ctx context.Context) error {
	if j == nil || j.manager == nil {
		return fmt.Errorf("job manager 未初始化")
	}
	return j.manager.Clear(ctx)
}

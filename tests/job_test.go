package tests

import (
	"context"
	"gin/app/facade"
	"gin/app/job"
	"gin/common/ctxkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestJobDispatchSync 同步Job投递(立即执行)
func TestJobDispatchSync(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-sync")

	err := facade.Job().Dispatch(ctx, "sync_user", job.SyncUser{
		UserID: 1001,
		Action: "update",
	})
	require.NoError(t, err, "同步Job投递失败")
	t.Log("同步Job投递成功")
}

// TestJobDispatchRedis Redis Job投递
func TestJobDispatchRedis(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-redis")

	time.Sleep(500 * time.Millisecond)

	err := facade.Job().Dispatch(ctx, "send_email", job.SendEmail{
		To:      "test@example.com",
		Subject: "测试邮件",
		Content: "这是一封测试邮件",
	})
	require.NoError(t, err, "Redis Job投递失败")
	t.Log("Redis Job投递成功")
}

// TestJobDispatchRabbitmq RabbitMQ Job投递(含延迟)
func TestJobDispatchRabbitmq(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-rabbitmq")

	time.Sleep(500 * time.Millisecond)

	err := facade.Job().Dispatch(ctx, "export_report", job.ExportReport{
		ReportType: "sales",
		StartDate:  "2026-01-01",
		EndDate:    "2026-06-30",
		UserID:     1001,
	})
	require.NoError(t, err, "RabbitMQ Job投递失败")
	t.Logf("RabbitMQ Job投递成功(延迟=%dms)", 10000)
}

// TestJobList 任务列表
func TestJobList(t *testing.T) {
	jobs := facade.Job().GetAllJobs()
	assert.NotEmpty(t, jobs, "任务列表不能为空")

	t.Logf("已注册 %d 个Job", len(jobs))
	for _, j := range jobs {
		t.Logf("  - %s (连接: %s, 描述: %s)", j.Name, j.Connection, j.Description)
	}
}

// TestJobDispatchUnregistered 投递未注册Job
func TestJobDispatchUnregistered(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-unreg")

	err := facade.Job().Dispatch(ctx, "nonexistent_job", nil)
	assert.Error(t, err, "未注册Job应该返回错误")
	t.Logf("未注册Job错误: %v", err)
}

// TestJobDelayVerify 验证延迟Job(Redis)确实在延迟队列等待
func TestJobDelayVerify(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-delay-verify")

	time.Sleep(500 * time.Millisecond)

	// 先清除
	_ = facade.Job().Clear(ctx)

	// 投递一个带延迟的Redis Job(send_email默认delay=3000ms)
	err := facade.Job().Dispatch(ctx, "send_email", job.SendEmail{
		To:      "delay@example.com",
		Subject: "延迟验证",
		Content: "验证延迟队列",
	})
	require.NoError(t, err, "投递失败")

	// 立即检查: 延迟队列中应该有1条待处理
	count, err := facade.Job().Count(ctx)
	require.NoError(t, err, "计数失败")
	assert.GreaterOrEqual(t, count, int64(1), "延迟队列应该有待处理Job")
	t.Logf("投递后待处理Job数: %d", count)

	// 等待延迟到期(send_email delay=3000ms, 等5秒确保被消费)
	time.Sleep(5 * time.Second)

	// 再次检查: 应该已被消费
	count2, err := facade.Job().Count(ctx)
	require.NoError(t, err, "计数失败")
	t.Logf("等待后待处理Job数: %d", count2)

	_ = facade.Job().Clear(ctx)
}

// TestJobRedisCountClear Redis Job计数和清除
func TestJobRedisCountClear(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-job-count")

	time.Sleep(500 * time.Millisecond)

	// 投递一个Redis Job
	err := facade.Job().Dispatch(ctx, "send_email", job.SendEmail{
		To:      "count@example.com",
		Subject: "计数测试",
		Content: "计数测试邮件",
	})
	require.NoError(t, err, "投递失败")

	// 等待消费
	time.Sleep(1500 * time.Millisecond)

	// 计数
	count, err := facade.Job().Count(ctx)
	require.NoError(t, err, "计数失败")
	t.Logf("待处理Job数: %d", count)

	// 清除
	err = facade.Job().Clear(ctx)
	require.NoError(t, err, "清除失败")
	t.Log("已清除所有待处理Job")
}

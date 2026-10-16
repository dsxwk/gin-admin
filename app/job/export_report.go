package job

import (
	"context"
	"fmt"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
)

// ExportReport 导出报表参数
type ExportReport struct {
	ReportType string `json:"reportType"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	UserID     int64  `json:"userId"`
}

// ExportReportJob 导出报表任务
type ExportReportJob struct{}

// NewExportReportJob 创建导出报表任务
func NewExportReportJob() *servicejob.Adapter[ExportReport] {
	return servicejob.NewAdapter(
		servicejob.Config{
			Name:        "export_report",
			Description: "导出报表任务",
			Connection:  "rabbitmq",
			Retry:       0,
			Delay:       10000,
		},
		(&ExportReportJob{}).Handle,
	)
}

// Handle 处理导出报表任务
func (j *ExportReportJob) Handle(_ context.Context, data *ExportReport) error {
	facade.Log().Info(fmt.Sprintf("Job [export_report] 导出报表: type=%s user=%d", data.ReportType, data.UserID))
	return nil
}

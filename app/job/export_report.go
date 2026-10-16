package job

import (
	"context"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
	"fmt"
)

// ExportReport 导出报表参数
type ExportReport struct {
	ReportType string `json:"reportType"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	UserID     int64  `json:"userId"`
}

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
		handleExportReport,
	)
}

// handleExportReport 处理导出报表任务
func handleExportReport(_ context.Context, data *ExportReport) error {
	facade.Log().Info(fmt.Sprintf("Job [export_report] 导出报表: type=%s user=%d", data.ReportType, data.UserID))
	return nil
}

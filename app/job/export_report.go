package job

import (
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
)

// ExportReportJob 导出报表任务
type ExportReportJob struct{}

// ExportReport 导出报表参数
type ExportReport struct {
	ReportType string `json:"report_type"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	UserID     int64  `json:"user_id"`
}

func (j *ExportReportJob) Name() string        { return "export_report" }
func (j *ExportReportJob) Description() string { return "导出报表任务" }
func (j *ExportReportJob) Connection() string  { return "rabbitmq" }
func (j *ExportReportJob) Retry() int          { return 0 }
func (j *ExportReportJob) Delay() int64        { return 10000 }
func (j *ExportReportJob) NewPayload() any     { return &ExportReport{} }

func (j *ExportReportJob) Handle(payload any) error {
	data := payload.(*ExportReport)
	facade.Log().Info(pkg.Sprintf("Job [export_report] 导出报表: type=%s user=%d", data.ReportType, data.UserID))
	return nil
}

func init() {
	job.Register(&ExportReportJob{})
}

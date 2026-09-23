package job

import servicejob "gin/pkg/serviceprovider/job"

// All 任务列表
func All() []servicejob.Job {
	return []servicejob.Job{
		&SendEmailJob{},
		&ExportReportJob{},
		&SyncUserJob{},
	}
}

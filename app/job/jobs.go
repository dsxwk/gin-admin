package job

import servicejob "gin/pkg/serviceprovider/job"

// Jobs 获取任务列表
func Jobs() []servicejob.Job {
	return []servicejob.Job{
		&SendEmailJob{},
		&ExportReportJob{},
		&SyncUserJob{},
	}
}

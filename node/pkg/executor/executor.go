package executor

import "go-job/internal/models"

type IExecutor interface {
	Run()                                            // 为了实现cron库的Job接口
	Execute() error                                  // 执行方法
	SetOnStatusChange(func(status models.JobStatus)) // 回调任务状态
}

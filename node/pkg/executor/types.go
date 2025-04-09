package executor

import "go-job/internal/model"

type IExecutor interface {
	Run()                                            // 为了实现cron库的Job接口
	Execute() (string, error)                        // 执行方法
	OnResultChange(func(result model.JobExecResult)) // 回调任务状态
}

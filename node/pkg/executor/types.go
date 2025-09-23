package executor

import "go-job/internal/model"

type IExecutor interface {
	Run()                                            // 为了实现cron库的Job接口
	Execute() (string, error)                        // 执行方法
	OnResultChange(func(result model.JobExecResult)) // 注册回调任务状态
	ResultCallback(output string, err error)         // 执行回调
	AfterExecute(err error)                          // 执行方法前的调用
	BeforeExecute()                                  // 执行方法后的调用
}

package executor

import (
	"go-job/internal/model"
	"log/slog"
)

type retryExecutor struct {
	executor IExecutor
	retries  int
}

func (r *retryExecutor) ResultCallback(output string, err error) {
	r.executor.ResultCallback(output, err)
}

func (r *retryExecutor) AfterExecute(err error) {
	r.executor.AfterExecute(err)
}

func (r *retryExecutor) BeforeExecute() {
	r.executor.BeforeExecute()
}

func NewRetryExecutor(executor IExecutor, retries int) IExecutor {
	return &retryExecutor{
		executor: executor,
		retries:  retries,
	}
}

func (r *retryExecutor) Run() {
	r.executor.BeforeExecute()
	output, err := r.Execute()
	r.executor.AfterExecute(err)
	r.executor.ResultCallback(output, err)
}

func (r *retryExecutor) Execute() (string, error) {
	var output string
	var err error
	for i := 0; i < r.retries; i++ {
		output, err = r.executor.Execute()
		if err == nil {
			return output, nil
		}
		slog.Error("executor failed", "attempt", i+1, "err", err)
	}
	slog.Info("retry over, exec failed", "retries", r.retries)
	return output, err
}

func (r *retryExecutor) OnResultChange(fn func(result model.JobExecResult)) {
	r.executor.OnResultChange(fn)
}

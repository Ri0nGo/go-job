package executor

import (
	"go-job/internal/model"
)

type ExampleStrategy struct {
	executor IExecutor
}

func (e *ExampleStrategy) Run() {
	e.executor.Run()
}

func (e *ExampleStrategy) Execute() (string, error) {
	// todo 这里可以做执行前的一些限制
	return e.executor.Execute()

}

func (e *ExampleStrategy) OnResultChange(f func(result model.JobExecResult)) {
	e.executor.OnResultChange(f)
}

func (e *ExampleStrategy) ResultCallback(output string, err error) {
	e.executor.ResultCallback(output, err)
}

func (e *ExampleStrategy) AfterExecute(err error) {
	e.executor.AfterExecute(err)
}

func (e *ExampleStrategy) BeforeExecute() {
	e.executor.BeforeExecute()
}

func NewExampleStrategy(executor IExecutor) IExecutor {
	return &ExampleStrategy{
		executor: executor,
	}

}

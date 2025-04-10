package executor

import (
	"go-job/internal/dto"
	"go-job/internal/model"
)

type ExecutorFactory func(req dto.ReqJob) IExecutor

var factories = make(map[model.ExecType]ExecutorFactory)

func Register(execType model.ExecType, factory ExecutorFactory) {
	factories[execType] = factory
}

func GetExecutor(execType model.ExecType) (ExecutorFactory, bool) {
	f, ok := factories[execType]
	return f, ok
}

func init() {
	Register(model.ExecTypeFile, func(req dto.ReqJob) IExecutor {
		return NewFileExecutor(req.Id, req.Name, req.Filename)
	})
}

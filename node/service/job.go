package service

import (
	"context"
	"errors"
	"go-job/internal/dto"
	"go-job/internal/models"
	"go-job/node/pkg/executor"
	"go-job/node/pkg/job"
)

type IJobService interface {
	AddJob(ctx context.Context, req dto.ReqJob) error
}

type JobService struct {
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) AddJob(ctx context.Context, req dto.ReqJob) error {
	// 获取对于的执行器
	exec, err := s.newExecutor(ctx, req)
	if err != nil {
		return err
	}

	// 构造job对象
	ctx, cancel := context.WithCancel(ctx)
	jj := job.NewJob(ctx, cancel, req, exec)
	if err = jj.ParseCrontab(); err != nil {
		return err
	}
	// 设置状态回调事件
	exec.SetOnStatusChange(jj.OnStatusChange)
	jj.Start()

	job.AddJob(jj)
	return nil
}

func (s *JobService) newExecutor(ctx context.Context, req dto.ReqJob) (executor.IExecutor, error) {
	var exec executor.IExecutor

	switch req.ExecType {
	case models.ExecTypeFile:
		exec = executor.NewFileExecutor(req.Filename)
	default:
		return nil, errors.New("invalid exec type")
	}
	return exec, nil
}

package service

import (
	"context"
	"errors"
	"go-job/internal/models"
	"go-job/node/pkg/executor"
	"go-job/node/pkg/job"
)

type IJobService interface {
	AddJob(ctx context.Context, job models.Job) error
}

type JobService struct {
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) AddJob(ctx context.Context, jobDAO models.Job) error {
	// 获取对于的执行器
	exec, err := s.newExecutor(ctx, jobDAO)
	if err != nil {
		return err
	}

	// 构造job对象
	ctx, cancel := context.WithCancel(ctx)
	jj := job.NewJob(ctx, cancel, jobDAO, exec)
	if err = jj.ParseCrontab(); err != nil {
		return err
	}
	jj.Start()

	job.AddJob(jj)
	return nil
}

func (s *JobService) newExecutor(ctx context.Context, job models.Job) (executor.IExecutor, error) {
	var exec executor.IExecutor

	switch job.ExecType {
	case models.ExecTypeFile:
		exec = executor.NewFileExecutor()
	default:
		return nil, errors.New("invalid exec type")
	}
	return exec, nil
}

package service

import (
	"context"
	"errors"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/node/pkg/executor"
	"go-job/node/pkg/job"
)

var (
	errJobNotFound = errors.New("job not found")
)

type IJobService interface {
	AddJob(ctx context.Context, req dto.ReqJob) error
	DeleteJob(ctx context.Context, id int)
	UpdateJob(ctx context.Context, req dto.ReqJob) error
	GetJob(ctx context.Context, id int) (*job.Job, error)
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
	if err = jj.BuildCrontab(); err != nil {
		return err
	}
	// 设置状态回调事件
	exec.SetOnStatusChange(jj.OnStatusChange)
	jj.Start()

	job.AddJob(jj)
	return nil
}

func (s *JobService) DeleteJob(ctx context.Context, id int) {
	if j, err := s.GetJob(ctx, id); err == nil {
		s.removeJob(j)
	}
}

func (s *JobService) removeJob(j *job.Job) {
	j.Cancel()
	j.Stop()
	job.RemoveJob(j.JobMeta.Id)
}

func (s *JobService) UpdateJob(ctx context.Context, req dto.ReqJob) error {
	j, ok := job.GetJob(req.Id)
	if !ok {
		return errJobNotFound
	}
	s.removeJob(j)

	if err := s.AddJob(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *JobService) GetJob(ctx context.Context, id int) (*job.Job, error) {
	j, ok := job.GetJob(id)
	if !ok {
		return nil, errJobNotFound
	}
	return j, nil
}

func (s *JobService) newExecutor(ctx context.Context, req dto.ReqJob) (executor.IExecutor, error) {
	var exec executor.IExecutor

	switch req.ExecType {
	case model.ExecTypeFile:
		exec = executor.NewFileExecutor(req.Filename)
	default:
		return nil, errors.New("invalid exec type")
	}
	return exec, nil
}

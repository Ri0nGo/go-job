package service

import (
	"context"
	"errors"
	"fmt"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/node/pkg/executor"
	"go-job/node/pkg/job"
	"log/slog"
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
	// job 已经存在，不能重复添加
	if _, ok := job.GetJob(req.Id); ok {
		slog.Error("job already exists, don't add job", "job id", req.Id,
			"job name", req.Name)
		return nil
	}

	jj, err := s.buildJobItem(ctx, req)
	if err != nil {
		return err
	}
	if req.Active == model.JobStart {
		jj.Start()
	}

	job.AddJob(jj)
	return nil
}

func (s *JobService) buildJobItem(ctx context.Context, req dto.ReqJob) (*job.Job, error) {
	// 获取对于的执行器
	exec, err := s.newExecutor(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构造job对象
	ctx, cancel := context.WithCancel(ctx)
	jj := job.NewJob(ctx, cancel, req, exec)
	if err = jj.BuildCrontab(); err != nil {
		return nil, err
	}
	// 设置状态回调事件
	exec.OnResultChange(jj.OnResultChange)
	return jj, nil
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

	jj, err := s.buildJobItem(ctx, req)
	if err != nil {
		return err
	}
	if req.Active == model.JobStart {
		jj.Start()
	}

	job.AddJob(jj)
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
	factory, ok := executor.GetExecutor(req.ExecType)
	if !ok {
		return nil, fmt.Errorf("unsupported executor type: %v", req.ExecType)
	}
	return factory(req), nil
}

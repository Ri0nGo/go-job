package service

import (
	"errors"
	"github.com/robfig/cron/v3"
	"go-job/internal/model"
	"go-job/internal/upload"
	"go-job/master/repo"
)

type IJobService interface {
	GetJob(id int) (model.Job, error)
	GetJobList(page model.Page) (model.Page, error)
	AddJob(job model.Job) error
	DeleteJob(id int) error
	UpdateJob(job model.Job) error
}

type JobService struct {
	JobRepo repo.IJobRepo
}

func (j *JobService) GetJob(id int) (model.Job, error) {
	return j.JobRepo.QueryById(id)
}

func (j *JobService) GetJobList(page model.Page) (model.Page, error) {
	return j.JobRepo.QueryList(page)
}

func (j *JobService) AddJob(job model.Job) error {
	if err := j.validExecType(job); err != nil {
		return err
	}

	if err := j.parseCrontab(job.CronExpr); err != nil {
		return err
	}
	return j.JobRepo.Inserts([]model.Job{job})
}

func (j *JobService) validExecType(job model.Job) error {
	switch job.ExecType {
	case model.ExecTypeFile:
		fileMeta, ok := upload.GetFileMeta(job.FileKey)
		if !ok {
			return errors.New("file not exist")
		}
		job.Internal.FileMeta = fileMeta
	default:
		return errors.New("not support exec type")
	}
	return nil
}

func (j *JobService) parseCrontab(cronExpr string) error {
	_, err := cron.ParseStandard(cronExpr)
	return err
}

func (j *JobService) DeleteJob(id int) error {
	return j.JobRepo.Delete(id)
}

func (j *JobService) UpdateJob(job model.Job) error {
	if err := j.validExecType(job); err != nil {
		return err
	}
	if err := j.parseCrontab(job.CronExpr); err != nil {
		return err
	}
	if _, err := j.GetJob(job.Id); err != nil {
		return err
	}
	return j.JobRepo.Update(job)
}

func NewJobService(jobRepo repo.IJobRepo) IJobService {
	return &JobService{
		JobRepo: jobRepo,
	}
}

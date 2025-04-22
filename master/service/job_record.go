package service

import (
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/notify"
	"go-job/master/repo"
)

type IJobRecordService interface {
	GetJobRecord(id int) (model.JobRecord, error)
	GetJobRecordList(page model.Page, jobId int) (model.Page, error)
	AddJobRecord(req model.CallbackJobResult) error
	DeleteJobRecord(id int) error
}

type JobRecordService struct {
	JobRecordRepo repo.IJobRecordRepo
	notifyStore   notify.INotifyStore
}

func (s *JobRecordService) GetJobRecord(id int) (model.JobRecord, error) {
	return s.JobRecordRepo.QueryById(id)
}

func (s *JobRecordService) GetJobRecordList(page model.Page, jobId int) (model.Page, error) {
	return s.JobRecordRepo.QueryList(page, jobId)
}

func (s *JobRecordService) AddJobRecord(req model.CallbackJobResult) error {
	jobRecord := model.JobRecord{
		JobId:        req.JobID,
		StartTime:    utils.TimestampToTime(req.StartTime),
		EndTime:      utils.TimestampToTime(req.EndTime),
		Status:       req.Status,
		NextExecTime: utils.TimestampToTime(req.NextExecTime),
		Duration:     req.Duration,
		Output:       req.Output,
		Error:        req.Error,
	}
	return s.JobRecordRepo.Insert(&jobRecord)
}

func (s *JobRecordService) DeleteJobRecord(id int) error {
	return s.JobRecordRepo.Delete(id)
}

func NewJobRecordService(nodeRepo repo.IJobRecordRepo, notify notify.INotifyStore) IJobRecordService {
	return &JobRecordService{
		JobRecordRepo: nodeRepo,
		notifyStore:   notify,
	}
}

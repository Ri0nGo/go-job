package service

import (
	"context"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/notify"
	"go-job/master/repo"
	"time"
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
	if jobId == 0 {
		if page.PageSize <= 0 {
			page.PageSize = 20
		} else if page.PageSize > 50 {
			page.PageSize = 50
		}
		return s.JobRecordRepo.QueryLastList(page)
	} else {
		return s.JobRecordRepo.QueryList(page, jobId)
	}
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
	if err := s.JobRecordRepo.Insert(&jobRecord); err != nil {
		return err
	}
	if nc, ok := s.notifyStore.Get(context.Background(), req.JobID); ok {
		s.notifyStore.PushNotifyUnit(context.Background(), req.JobID, s.jobToNotifyUnit(req, nc))
	}
	return nil
}

func (s *JobRecordService) jobToNotifyUnit(req model.CallbackJobResult, nc notify.NotifyConfig) notify.NotifyUnit {
	return notify.NotifyUnit{
		NotifyConfig: notify.NotifyConfig{
			JobID:          nc.JobID,
			Name:           nc.Name,
			NotifyStrategy: nc.NotifyStrategy,
			NotifyType:     nc.NotifyType,
			NotifyMark:     nc.NotifyMark,
		},
		StartExecTime: time.Unix(req.StartTime, 0).Local(),
		Status:        req.Status,
		Duration:      req.Duration,
		Output:        req.Output,
		Error:         req.Error,
	}
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

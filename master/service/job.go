package service

import "go-job/master/repo"

type IJobService interface {
}

type JobService struct {
	JobRepo repo.IJobRepo
}

func NewJobService(jobRepo repo.IJobRepo) IJobService {
	return &JobService{
		JobRepo: jobRepo,
	}
}

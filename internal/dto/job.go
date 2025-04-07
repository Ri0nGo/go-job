package dto

import "go-job/internal/models"

// ----------- node ----------- //

type ReqJob struct {
	models.Job
	Filename string `json:"filename"`
}

type ReqId struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type RespJob struct {
	Id            int              `json:"id"`
	Name          string           `json:"name"`
	ExecType      models.ExecType  `json:"exec_type"`
	RunningStatus models.JobStatus `json:"running_status"`
	FileName      string           `json:"file_name"`
}

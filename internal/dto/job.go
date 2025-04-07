package dto

import "go-job/internal/model"

// ----------- node ----------- //

type ReqJob struct {
	model.Job
	Filename string `json:"filename"`
}

type ReqId struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type RespJob struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	ExecType      model.ExecType  `json:"exec_type"`
	RunningStatus model.JobStatus `json:"running_status"`
	FileName      string          `json:"file_name"`
}

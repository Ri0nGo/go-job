package dto

import (
	"go-job/internal/model"
)

// ----------- node ----------- //

type ReqJob struct {
	Id       int                 `json:"id"`
	Name     string              `json:"name" binding:"required"`       // 任务名称
	ExecType model.ExecType      `json:"exec_type"  binding:"required"` // 任务类型
	CronExpr string              `json:"cron_expr" binding:"required"`  // crontab 表达式
	Active   model.JobActiveType `json:"active" binding:"required"`
	Filename string              `json:"filename"`
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

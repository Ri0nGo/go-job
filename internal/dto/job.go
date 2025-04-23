package dto

import (
	"go-job/internal/model"
	"time"
)

// ----------- Node DTO ----------- //

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

type ReqJobRecords struct {
	model.Page
	JobId int `json:"job_id" form:"job_id" binding:"required"`
}

type RespNodeJob struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	ExecType      model.ExecType  `json:"exec_type"`
	RunningStatus model.JobStatus `json:"running_status"`
	FileName      string          `json:"filename"`
}

// ----------- Master DTO ----------- //

type RespJob struct {
	Id            int                 `json:"id"`
	NodeID        int                 `json:"node_id"`
	Name          string              `json:"name"`
	ExecType      model.ExecType      `json:"exec_type"`
	CronExpr      string              `json:"cron_expr"`
	NodeName      string              `json:"node_name"`
	Active        model.JobActiveType `json:"active"`
	HasPermission bool                `json:"has_permission"`
	FileName      string              `json:"filename"`
	CreatedTime   time.Time           `json:"created_time"`
}

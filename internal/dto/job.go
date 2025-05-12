package dto

import (
	"go-job/internal/model"
	"time"
)

// ----------- Node DTO ----------- //

// 发送数据到node的struct
type ReqNodeJob struct {
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
	JobId int `json:"job_id" form:"job_id"`
}

type RespNodeJob struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	ExecType      model.ExecType  `json:"exec_type"`
	RunningStatus model.JobStatus `json:"running_status"`
	FileName      string          `json:"filename"`
}

// ----------- Master DTO ----------- //

type ReqJob struct {
	Id             int                  `json:"id"`
	Name           string               `json:"name"`                         // 任务名称
	ExecType       model.ExecType       `json:"exec_type" binding:"required"` // 任务类型
	CronExpr       string               `json:"cron_expr" binding:"required"` // crontab 表达式
	CreatedTime    time.Time            `json:"created_time"`
	UpdatedTime    time.Time            `json:"updated_time"`
	Active         model.JobActiveType  `json:"active"`
	NodeID         int                  `json:"node_id"`
	NotifyStatus   model.NotifyStatus   `json:"notify_status"`   // 通知启停
	NotifyType     model.NotifyType     `json:"notify_type" `    // 通知类型，邮件，短信等
	NotifyStrategy model.NotifyStrategy `json:"notify_strategy"` // 通知策略
	NotifyMark     string               `json:"notify_mark" `    // 通知方式的具体内容，可能是邮箱地址，可能是外链。
	UserId         int                  `json:"user_id"`
	FileName       string               `json:"filename" ` // 文件名
	FileKey        string               `json:"file_key"`  // 文件key
}

type RespJob struct {
	Id             int                  `json:"id"`
	NodeID         int                  `json:"node_id"`
	Name           string               `json:"name"`
	ExecType       model.ExecType       `json:"exec_type"`
	CronExpr       string               `json:"cron_expr"`
	NodeName       string               `json:"node_name"`
	Active         model.JobActiveType  `json:"active"`
	FileName       string               `json:"filename"`
	CreatedTime    time.Time            `json:"created_time"`
	NotifyStatus   model.NotifyStatus   `json:"notify_status"`   // 通知启停
	NotifyType     model.NotifyType     `json:"notify_type"`     // 通知类型，邮件，短信等
	NotifyStrategy model.NotifyStrategy `json:"notify_strategy"` // 通知策略
	NotifyMark     string               `json:"notify_mark"`     // 通知方式的具体内容，可能是邮箱地址，可能是外链。
	UserId         int                  `json:"user_id"`
}

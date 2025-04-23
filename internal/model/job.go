package model

import (
	"go-job/internal/upload"
	"strconv"
	"time"
)

type JobStatus int

const (
	Pending JobStatus = iota
	Running
	Success
	Failed
)

func (s JobStatus) String() string {
	switch s {
	case Pending:
		return "等待中"
	case Running:
		return "运行中"
	case Success:
		return "成功"
	case Failed:
		return "失败"
	default:
		return strconv.Itoa(int(s))
	}
}

type JobActiveType int8

const (
	JobStart JobActiveType = iota + 1
	JobStop
)

type Job struct {
	Id             int            `json:"id" gorm:"primary_key"`
	Name           string         `json:"name" binding:"required"`                              // 任务名称
	ExecType       ExecType       `json:"exec_type" gorm:"column:exec_type" binding:"required"` // 任务类型
	CronExpr       string         `json:"cron_expr" gorm:"column:cron_expr" binding:"required"` // crontab 表达式
	CreatedTime    time.Time      `json:"created_time" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime    time.Time      `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"`
	Active         JobActiveType  `json:"active" gorm:"column:active;default:1"`
	Internal       Internal       `gorm:"serializer:json;column:internal"`
	NodeID         int            `json:"node_id" gorm:"column:node_id" binding:"required"`
	NotifyStatus   NotifyStatus   `json:"notify_status" gorm:"column:notify_status;default:1"` // 通知启停
	NotifyType     NotifyType     `json:"notify_type" gorm:"column:notify_type"`               // 通知类型，邮件，短信等
	NotifyStrategy NotifyStrategy `json:"notify_strategy" gorm:"column:notify_strategy"`       // 通知策略
	NotifyMark     string         `json:"notify_mark" gorm:"column:notify_mark"`               // 通知方式的具体内容，可能是邮箱地址，可能是外链。
	UserId         int            `json:"user_id" gorm:"column:user_id"`
	HasPermission  bool           `json:"has_permission" gorm:"-"`
	FileName       string         `json:"filename" gorm:"-"` // 文件名
	FileKey        string         `json:"file_key" gorm:"-"` // 文件key
}

func (Job) TableName() string {
	return "job"
}

type Internal struct {
	FileMeta upload.FileMeta `json:"file_meta"`
}

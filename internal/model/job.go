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
		return "pending"
	case Running:
		return "running"
	case Success:
		return "success"
	case Failed:
		return "failed"
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
	Id          int           `json:"id" gorm:"primary_key"`
	Name        string        `json:"name" binding:"required"`                              // 任务名称
	ExecType    ExecType      `json:"exec_type" gorm:"column:exec_type" binding:"required"` // 任务类型
	CronExpr    string        `json:"cron_expr" gorm:"column:cron_expr" binding:"required"` // crontab 表达式
	CreatedTime time.Time     `json:"created_time" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time     `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"`
	Active      JobActiveType `json:"active" gorm:"column:active;default:1"`
	Internal    Internal      `gorm:"serializer:json;column:internal"`
	NodeID      int           `json:"node_id" gorm:"column:node_id" binding:"required"`
	FileName    string        `json:"filename" gorm:"-"` // 文件名
	FileKey     string        `json:"file_key" gorm:"-"` // 文件key
}

func (Job) TableName() string {
	return "job"
}

type Internal struct {
	FileMeta upload.FileMeta `json:"file_meta"`
}

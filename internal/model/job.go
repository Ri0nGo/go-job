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

type Job struct {
	Id        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`                              // 任务名称
	ExecType  ExecType  `json:"exec_type" gorm:"column:exec_type"` // 任务类型
	CronExpr  string    `json:"cron_expr" gorm:"column:cron_expr"` // crontab 表达式
	CreatedAt time.Time `json:"created_at" gorm:"column:created_time;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_time;autoUpdateTime"`
	Internal  Internal  `gorm:"serializer:json;column:internal"`
	FileName  string    `json:"filename" gorm:"-"` // 文件名
	FileKey   string    `json:"file_key" gorm:"-"` // 文件key
}

func (Job) TableName() string {
	return "job"
}

type Internal struct {
	FileMeta upload.FileMeta `json:"file_meta"`
}

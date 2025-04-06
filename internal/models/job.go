package models

import (
	"go-job/internal/upload"
	"strconv"
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
	Id       int      `json:"id"`
	Name     string   `json:"name"`      // 任务名称
	ExecType ExecType `json:"exec_type"` // 任务类型
	Crontab  string   `json:"crontab"`   // crontab 表达式
	Tags     []string `json:"tags"`      // 标签
	Internal Internal `json:"internal"`
}

type Internal struct {
	FileMeta upload.FileMeta `json:"file_meta"`
}

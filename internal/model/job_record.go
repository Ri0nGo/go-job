package model

import "time"

type JobExecResult struct {
	StartTime int64     `json:"start_time"`
	EndTime   int64     `json:"end_time"`
	Duration  float64   `json:"duration"`
	Status    JobStatus `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
}

type CallbackJobResult struct {
	JobExecResult
	JobID        int   `json:"job_id"` // job id
	NextExecTime int64 `json:"next_exec_time"`
}

type JobRecord struct {
	Id           int       `json:"id"`
	JobId        int       `json:"job_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	NextExecTime time.Time `json:"next_exec_time"`
	Duration     float64   `json:"duration"`
	Status       JobStatus `json:"status"`
	Output       string    `json:"output"`
	Error        string    `json:"error"`
}

type JobLastRecord struct {
	Id        int       `json:"id"`
	JobId     int       `json:"job_id"`
	NodeId    int       `json:"node_id"`
	JobName   string    `json:"job_name"`
	NodeName  string    `json:"node_name"`
	Status    JobStatus `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// JobRecordSummary jobRecord 简化结构体，移除了大范围输出
type JobRecordSummary struct {
	Id           int       `json:"id"`
	JobId        int       `json:"job_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	NextExecTime time.Time `json:"next_exec_time"`
	Duration     float64   `json:"duration"`
	Status       JobStatus `json:"status"`
}

type JobRecordDayStatusCount struct {
	Date   time.Time `json:"date"`
	Status JobStatus `json:"status" gorm:"column:status"`
	Count  int       `json:"count" gorm:"column:count"`
}

type JobRecordJobStatusCount struct {
	JobId  int       `json:"job_id" gorm:"column:job_id"`
	Status JobStatus `json:"status" gorm:"column:status"`
	Count  int       `json:"count" gorm:"column:count"`
}

func (j *JobRecord) TableName() string {
	return "job_record"
}

func (j *JobRecordSummary) TableName() string {
	return "job_record"
}

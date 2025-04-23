package notify

import (
	"context"
	"go-job/internal/model"
	"time"
)

type NotifyConfig struct {
	JobID          int
	Name           string
	NotifyStrategy model.NotifyStrategy
	NotifyType     model.NotifyType
	NotifyMark     string
}

type INotifyStore interface {
	Set(ctx context.Context, jobId int, config NotifyConfig) error
	Get(ctx context.Context, jobId int) (nc NotifyConfig, ok bool)
	Delete(ctx context.Context, jobId int) error

	PushNotifyUnit(ctx context.Context, jobId int, unit NotifyUnit) error
}

type NotifyUnit struct {
	NotifyConfig
	Status        model.JobStatus // 执行状态
	StartExecTime time.Time
	Duration      float64 // 耗时
	Output        string
	Error         string
}

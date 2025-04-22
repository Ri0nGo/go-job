package notify

import (
	"context"
	"go-job/internal/model"
)

type NotifyConfig struct {
	JobID          int
	Name           string
	NotifyStrategy model.NotifyStrategy
	NotifyType     model.NotifyType
}

type INotifyStore interface {
	Set(ctx context.Context, jobId int, config NotifyConfig) error
	Get(ctx context.Context, jobId int) (nc NotifyConfig, ok bool)
	Delete(ctx context.Context, jobId int) error

	PushNotifyUnit(ctx context.Context, jobId int, unit NotifyUnit) error
}

type NotifyUnit struct {
	NotifyConfig
	Status   model.JobStatus // 执行状态
	Duration float64         // 耗时
	Output   string
	Error    string

	Email string // 邮箱
}

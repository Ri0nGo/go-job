package job

import (
	"context"
	"github.com/robfig/cron/v3"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/node/pkg/executor"
	"sync"
)

type JobMeta struct {
	Id       int            `json:"id"`
	Name     string         `json:"name"`      // 任务名称
	ExecType model.ExecType `json:"exec_type"` // 任务类型
	CronExpr string         `json:"cron_expr"` // crontab 表达式
	FileName string         `json:"file_name"` // 本地存储的文件名
}

type Job struct {
	Ctx           context.Context
	Cancel        context.CancelFunc
	JobMeta       JobMeta
	Executor      executor.IExecutor
	Cron          *cron.Cron
	CronEntryID   cron.EntryID
	RunningStatus model.JobStatus
}

// JobManager 全局job管理
type JobManager struct {
	mux  sync.RWMutex
	jobs map[int]*Job
}

var jm = &JobManager{
	jobs: make(map[int]*Job),
}

func AddJob(job *Job) {
	jm.mux.Lock()
	defer jm.mux.Unlock()
	jm.jobs[job.JobMeta.Id] = job
}

func RemoveJob(id int) {
	jm.mux.Lock()
	defer jm.mux.Unlock()
	delete(jm.jobs, id)
}

func GetJob(id int) (*Job, bool) {
	jm.mux.RLock()
	defer jm.mux.RUnlock()
	job, ok := jm.jobs[id]
	return job, ok
}

func GetAllJobs() []*Job {
	jm.mux.RLock()
	defer jm.mux.RUnlock()
	jobs := make([]*Job, 0, len(jm.jobs))
	for _, job := range jm.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func NewJob(ctx context.Context, cancel context.CancelFunc, req dto.ReqJob, iExecutor executor.IExecutor) *Job {
	return &Job{
		Ctx:    ctx,
		Cancel: cancel,
		JobMeta: JobMeta{
			Id:       req.Id,
			Name:     req.Name,
			ExecType: req.ExecType,
			CronExpr: req.CronExpr,
			FileName: req.Filename,
		},
		Executor:      iExecutor,
		RunningStatus: model.Pending,
	}
}

func (j *Job) BuildCrontab() error {
	c := cron.New(cron.WithSeconds())
	entryID, err := c.AddJob(j.JobMeta.CronExpr, j.Executor)
	if err != nil {
		return err
	}
	j.Cron = c
	j.CronEntryID = entryID
	return nil
}

func (j *Job) OnStatusChange(status model.JobStatus) {
	j.RunningStatus = status
}

func (j *Job) Start() {
	j.Cron.Start()
}

func (j *Job) Stop() context.Context {
	return j.Cron.Stop()
}

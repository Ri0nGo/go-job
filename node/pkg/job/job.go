package job

import (
	"context"
	"github.com/robfig/cron/v3"
	"go-job/internal/models"
	"go-job/node/pkg/executor"
	"sync"
)

type Job struct {
	Ctx           context.Context
	Cancel        context.CancelFunc
	JobDAO        models.Job
	Executor      executor.IExecutor
	Cron          *cron.Cron
	CronEntryID   cron.EntryID
	RunningStatus JobStatus
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
	jm.jobs[job.JobDAO.Id] = job
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

func NewJob(ctx context.Context, cancel context.CancelFunc, job models.Job, iExecutor executor.IExecutor) *Job {
	return &Job{
		Ctx:           ctx,
		Cancel:        cancel,
		JobDAO:        job,
		Executor:      iExecutor,
		RunningStatus: Pending,
	}
}

func (j *Job) ParseCrontab() error {
	c := cron.New(cron.WithSeconds())
	entryID, err := c.AddJob(j.JobDAO.Crontab, j.Executor)
	if err != nil {
		return err
	}
	j.Cron = c
	j.CronEntryID = entryID
	return nil
}

func (j *Job) Start() {
	j.Cron.Start()
}

func (j *Job) Stop() {
	j.Cron.Stop()
}

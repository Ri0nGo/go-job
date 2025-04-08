package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go-job/internal/model"
	"go-job/internal/pkg/httpClient"
	"go-job/internal/upload"
	"go-job/master/pkg/paths"
	"go-job/master/repo"
	"log/slog"
	"time"
)

type IJobService interface {
	GetJob(id int) (model.Job, error)
	GetJobList(page model.Page) (model.Page, error)
	AddJob(job model.Job) error
	DeleteJob(id int) error
	UpdateJob(job model.Job) error
}

type JobService struct {
	JobRepo  repo.IJobRepo
	NodeRepo repo.INodeRepo
}

func (j *JobService) GetJob(id int) (model.Job, error) {
	return j.JobRepo.QueryById(id)
}

func (j *JobService) GetJobList(page model.Page) (model.Page, error) {
	return j.JobRepo.QueryList(page)
}

func (j *JobService) AddJob(job model.Job) error {
	if err := j.parseExecType(&job); err != nil {
		return err
	}

	if err := j.parseCrontab(job.CronExpr); err != nil {
		return err
	}
	err := j.JobRepo.Inserts([]model.Job{job})
	if err != nil {
		return err
	}

	node, err := j.NodeRepo.QueryById(job.NodeID)
	if err != nil {
		return errors.New("node not found")
	}
	j.sendJobToNode(node, job)
	return nil
}

func (j *JobService) sendJobToNode(node model.Node, job model.Job) {
	go func() {
		type reqJobNode struct {
			Id       int            `json:"id"`
			Name     string         `json:"name"`
			ExecType model.ExecType `json:"exec_type"`
			CronExpr string         `json:"cron_expr"`
			Filename string         `json:"filename"`
		}
		req := reqJobNode{
			Id:       job.Id,
			Name:     job.Name,
			ExecType: job.ExecType,
			CronExpr: job.CronExpr,
			Filename: job.Internal.FileMeta.UUIDFileName,
		}
		url := fmt.Sprintf("http://%s%s%s", node.Address,
			paths.NodeJobAPI.BasePath, paths.NodeJobAPI.Create)
		resp, err := httpClient.PostJson(context.Background(), url, req, 3*time.Second)
		if err != nil {
			slog.Error("send job to node error", "err", err)
		}
		nodeResp, err := httpClient.ParseResponse(resp)
		if err != nil {
			slog.Error("send job to node error", "err", err)
		}
		if nodeResp.Code != 0 {
			slog.Error("resp code isn't zero", "resp", resp)
		}
	}()
}

// validExecType 创建任务的时候检测
func (j *JobService) parseExecType(job *model.Job) error {
	switch job.ExecType {
	case model.ExecTypeFile:
		fileMeta, ok := upload.GetFileMeta(job.FileKey)
		if !ok {
			return errors.New("file not exist")
		}
		job.Internal.FileMeta = fileMeta
		upload.DeleteFileMeta(job.FileKey)
	default:
		return errors.New("not support exec type")
	}
	return nil
}

func (j *JobService) parseCrontab(cronExpr string) error {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronExpr)
	return err
}

func (j *JobService) DeleteJob(id int) error {
	return j.JobRepo.Delete(id)
}

func (j *JobService) UpdateJob(job model.Job) error {
	if err := j.parseCrontab(job.CronExpr); err != nil {
		return err
	}
	dbJob, err := j.GetJob(job.Id)
	if err != nil {
		return err
	}
	switch job.ExecType {
	case model.ExecTypeFile:
		if len(job.FileName) == 0 && len(job.FileKey) == 0 {
			return errors.New("file name or file key is empty")
		}
		// 执行文件修改了
		if len(job.FileKey) > 0 {
			fileMeta, b := upload.GetFileMeta(job.FileKey)
			if !b {
				return errors.New("file not exist")
			}
			job.Internal = dbJob.Internal
			job.Internal.FileMeta = fileMeta
			upload.DeleteFileMeta(job.FileKey)
		}
	default:
		return errors.New("not support exec type")
	}
	return j.JobRepo.Update(job)
}

func NewJobService(jobRepo repo.IJobRepo, nodeRepo repo.INodeRepo) IJobService {
	return &JobService{
		JobRepo:  jobRepo,
		NodeRepo: nodeRepo,
	}
}

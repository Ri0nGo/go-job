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
	"go-job/node/pkg/config"
	"log/slog"
	"os"
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

	node, err := j.NodeRepo.QueryById(job.NodeID)
	if err != nil {
		return errors.New("node not found")
	}

	// 先发送任务到节点中，若发送失败则直接返回
	err = j.sendDataToNode(job, node)
	if err != nil {
		return err
	}

	// 将数据插入数据库
	err = j.JobRepo.Inserts([]model.Job{job})
	if err != nil {
		// 若数据插入失败，则移除问题
		// TODO 移除任务
		err = j.removeJobInNode(job.Id)
		if err != nil {
			slog.Error("remove job failed in node", err, err.Error())
		}
		return err
	}

	return nil
}

func (j *JobService) sendDataToNode(job model.Job, node model.Node) error {
	switch job.ExecType {
	case model.ExecTypeFile:
		// 发送文件到节点
		err := j.sendJobFileInNode(job, node)
		if err != nil {
			return err
		}
		// 添加任务到节点
		err = j.sendJobInNode(job, node)
		if err != nil {
			return err
		}
	default:
		return errors.New("not support exec type")

	}
	return nil
}

// sendJobFileToNode 发送job文件到节点
func (j *JobService) sendJobFileInNode(job model.Job, node model.Node) error {
	fileColName := "file"
	filePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir,
		job.Internal.FileMeta.UUIDFileName)
	formData := map[string]string{
		"filename": job.Internal.FileMeta.UUIDFileName,
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s%s%s", node.Address,
		paths.NodeJobAPI.BasePath, paths.NodeJobAPI.Upload)
	resp, err := httpClient.PostFormDataWithFile(context.Background(), fileColName,
		f, url, formData, httpClient.DefaultTimeout)
	if err != nil {
		return err
	}
	nodeResp, err := httpClient.ParseResponse(resp)
	if err != nil {
		slog.Error("send job to node error", "err", err)
		return err
	}
	if nodeResp.Code != 0 {
		slog.Error("resp code isn't zero", "resp", resp)
		return errors.New("resp code isn't zero in send job file")
	}
	return nil
}

// sendJobToNode 发送任务到节点
func (j *JobService) sendJobInNode(job model.Job, node model.Node) error {
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

	resp, err := httpClient.PostJson(context.Background(), url, req, httpClient.DefaultTimeout)
	if err != nil {
		slog.Error("send job to node error", "err", err)
		return err
	}
	nodeResp, err := httpClient.ParseResponse(resp)
	if err != nil {
		slog.Error("send job to node error", "err", err)
		return err
	}
	if nodeResp.Code != 0 {
		slog.Error("resp code isn't zero", "resp", resp)
		return errors.New("resp code isn't zero in send job data")
	}
	return nil
}

// removeJobInNode 移除任务
func (j *JobService) removeJobInNode(id int) error {
	return nil

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

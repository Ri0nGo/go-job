package startup

import (
	"context"
	"fmt"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/httpClient"
	"go-job/internal/pkg/paths"
	"go-job/node/pkg/auth"
	"go-job/node/pkg/config"
	"go-job/node/service"
	"log/slog"
	"strconv"
	"time"
)

// 内层 data
type respJobListData struct {
	Total    int         `json:"total"`
	PageSize int         `json:"page_size"`
	PageNum  int         `json:"page_num"`
	Data     []model.Job `json:"data"`
	Order    string      `json:"order"`
	Sort     string      `json:"sort"`
}

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func SyncJobFromMaster(jobSvc service.IJobService) error {
	if err := auth.RefreshToken(); err != nil {
		return err
	}

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": auth.GetJwtToken(),
	}
	url := fmt.Sprintf("http://%s%s", config.App.Master.Address, paths.JobListAPI)
	params := map[string]string{
		"active":    strconv.Itoa(int(model.JobStart)),
		"page_num":  "1",
		"page_size": "9999",
	}
	resp, err := httpClient.GetJson(context.Background(), url, header, params, 3*time.Second)
	if err != nil {
		return err
	}
	parseResp, err := httpClient.ParseResponseWith[Resp[respJobListData]](resp)
	if err != nil {
		return err
	}

	for _, job := range parseResp.Data.Data {
		if err = jobSvc.AddJob(context.Background(), dto.ReqNodeJob{
			Id:       job.Id,
			Name:     job.Name,
			ExecType: job.ExecType,
			CronExpr: job.CronExpr,
			Active:   job.Active,
			Filename: job.UUIDFileName,
		}); err != nil {
			slog.Error("sync job failed", "id", job.Id, "name", job.Name, "err", err)
		}
	}
	return nil
}

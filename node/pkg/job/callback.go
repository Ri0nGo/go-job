package job

import (
	"context"
	"fmt"
	"go-job/internal/model"
	"go-job/internal/pkg/httpClient"
	"go-job/internal/pkg/paths"
	"go-job/node/pkg/auth"
	"go-job/node/pkg/config"
	"log/slog"
)

// CallbackJobResult 回传结果到master,
func CallbackJobResult(result model.CallbackJobResult) {
	err := auth.RefreshToken()
	if err != nil {
		slog.Error("refresh token err:", "err", err)
		return
	}
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": auth.GetJwtToken(),
	}
	url := fmt.Sprintf("http://%s%s", config.App.Master.Address, paths.JobRecordCreateAPI)
	resp, err := httpClient.PostJson(context.Background(), url, header, result, httpClient.DefaultTimeout)
	if err != nil {
		slog.Error("callback result error", err, "url", url, "resp", resp, "err", err)
		return
	}
	parseContent, err := httpClient.ParseResponse(resp)
	if err != nil {
		slog.Error("callback result parse error", "url", url, "resp", resp,
			"status", resp.Status(), "err", err)
		return
	}
	if parseContent.Code != 0 {
		slog.Error("code isn't zero in callback result", "parse content", parseContent)
		return
	}
}

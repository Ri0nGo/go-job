package httpClient

import (
	"context"
	"encoding/json"
	"errors"
	"go-job/internal/dto"
	"resty.dev/v3"
	"time"
)

var defaultRestyClient = resty.New()

func PostJson(ctx context.Context, url string, body any, timeout time.Duration) (*resty.Response, error) {
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

func PostJsonWithAuth(ctx context.Context, url string, body any, timeout time.Duration, auth string) (*resty.Response, error) {
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(auth).
		SetBody(body).
		Post(url)
}

func ParseResponse(resp *resty.Response) (dto.Response, error) {
	var commResp = dto.Response{}
	if resp == nil {
		return commResp, errors.New("resp is nil")
	}
	err := json.Unmarshal(resp.Bytes(), &commResp)
	if err != nil {
		return commResp, err
	}
	return commResp, nil
}

package httpClient

import (
	"context"
	"encoding/json"
	"errors"
	"go-job/internal/dto"
	"os"
	"path/filepath"
	"resty.dev/v3"
	"time"
)

var (
	defaultRestyClient = resty.New()
	DefaultTimeout     = time.Second * 3
)

func GetJson(ctx context.Context, url string, headers map[string]string, params map[string]string, timeout time.Duration) (*resty.Response, error) {
	jsonHeader(headers)
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeaders(headers).
		SetQueryParams(params).
		Get(url)
}

func PostJson(ctx context.Context, url string, headers map[string]string, body any, timeout time.Duration) (*resty.Response, error) {
	jsonHeader(headers)
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Post(url)
}

func jsonHeader(header map[string]string) map[string]string {
	if header == nil {
		header = map[string]string{
			"Content-Type": "application/json",
		}
	}
	return header
}
func PutJson(ctx context.Context, url string, headers map[string]string, body any, timeout time.Duration) (*resty.Response, error) {
	jsonHeader(headers)
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Put(url)
}

func Delete(ctx context.Context, url string, headers map[string]string, timeout time.Duration, params map[string]string) (*resty.Response, error) {
	jsonHeader(headers)
	c := defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx)
	if params != nil {
		c.SetQueryParams(params)
	}
	return c.Delete(url)
}

// PostFormDataWithFile 发送含文件的form-data
/*
	fileColName：文件字段名
	f：文件打开的os.File类型
*/
func PostFormDataWithFile(ctx context.Context, fileColName string, f *os.File, url string,
	formData map[string]string, timeout time.Duration) (*resty.Response, error) {
	filename := filepath.Base(f.Name())
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetFileReader(fileColName, filename, f).
		SetFormData(formData).
		Post(url)
}

//if _, err := os.Stat(filePath); os.IsNotExist(err) {
//	fmt.Println("文件不存在:", filePath)
//	return
//}
//
//// 获取文件名
//filename := filepath.Base(filePath)
//
//// 创建 Resty 客户端
//client := resty.New()
//
//// 发送 POST 请求
//resp, err := client.R().
//	SetFileReader("file", filename, mustOpen(filePath)). // 添加 file 类型字段
//	SetFormData(map[string]string{                       // 添加 string 类型字段
//		"filename": filename,
//	}).
//	Post("your-upload-url") // 将其替换为你的上传 API URL
//}

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

func ParseResponseWith[T any](resp *resty.Response) (*T, error) {
	if resp == nil {
		return nil, errors.New("resp is nil")
	}

	var commResp T
	if err := json.Unmarshal(resp.Bytes(), &commResp); err != nil {
		return nil, err
	}
	return &commResp, nil
}

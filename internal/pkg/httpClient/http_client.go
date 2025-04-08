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

func PostJson(ctx context.Context, url string, body any, timeout time.Duration) (*resty.Response, error) {
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

func PostJsonWithAuth(ctx context.Context, url string, body any,
	timeout time.Duration, auth string) (*resty.Response, error) {
	return defaultRestyClient.SetTimeout(timeout).R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(auth).
		SetBody(body).
		Post(url)
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

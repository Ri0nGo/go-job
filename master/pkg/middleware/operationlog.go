package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-job/internal/model"
	"go-job/master/database"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

var bufferSize = 1024

const (
	OperationDescAddJob              = "新增任务"
	OperationDescDeleteJob           = "删除任务"
	OperationDescUpdateJob           = "更新任务"
	OperationDescUploadFile          = "上传代码文件"
	OperationDescDownloadFile        = "下载代码文件"
	OperationDescAddNode             = "新增节点"
	OperationDescDeleteNode          = "删除节点"
	OperationDescUpdateNode          = "更新节点"
	OperationDescNodeInstallRef      = "新增依赖包"
	OperationDescAddUser             = "新增用户"
	OperationDescDeleteUser          = "删除用户"
	OperationDescUpdateUser          = "更新用户"
	OperationDescLogin               = "登录"
	OperationDescOAuth2Bind          = "绑定第三方账号"
	OperationDescOAuth2UnBind        = "解绑第三方账号"
	OperationDescSendEmailCode       = "发生邮箱验证码"
	OperationDescBindEmail           = "绑定邮箱"
	OperationDescOAuth2Login         = "第三方账号登录"
	OperationDescOAuth2GithubAuthURL = "Github开启授权"
	OperationDescOAuth2QQAuthURL     = "QQ开启授权"
)

func OperationLog(optDesc string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body   []byte
			userId int
		)

		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				slog.Error("read body err:", "err", err)
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			values, _ := url.ParseQuery(query)
			m := map[string]string{}
			for k, v := range values {
				if len(v) > 0 {
					m[k] = v[0]
				}
			}
			body, _ = json.Marshal(&m)
		}

		if strings.Contains(c.GetHeader("content-type"), "multipart/form-data") {
			file, err := c.FormFile("file") // "file" 默认获取file
			if err != nil {
				body = []byte("获取文件信息失败")
			} else {
				body = []byte(fmt.Sprintf("文件: %s (%d kb)", file.Filename, file.Size/1024))
			}
		} else {
			if len(body) > bufferSize {
				body = body[:bufferSize]
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		uc, err := GetUserClaim(c)
		if err != nil {
			slog.Error("get user claim err:", "err", err)
		} else {
			userId = uc.Uid
		}

		optLog := model.SystemOperationLog{
			UserId:     userId,
			ClientIP:   c.ClientIP(),
			UA:         c.Request.UserAgent(),
			Method:     c.Request.Method,
			URL:        c.Request.URL.Path,
			Request:    string(body),
			Title:      optDesc,
			StatusCode: c.Writer.Status(),
			Response:   writer.body.String(),
		}

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
			if len(optLog.Response) > bufferSize {
				optLog.Response = "超出记录长度"
			}
		}

		if err = database.CreateSysOptLogWithMySQL(optLog); err != nil {
			slog.Error("create sys opt log err:", "err", err)
		}

	}
}

// GetUserClaim 获取用户的UC信息
func GetUserClaim(ctx *gin.Context) (*model.UserClaims, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("user not exists in context")
	}
	uc, ok := value.(*model.UserClaims)
	if !ok {
		return nil, errors.New("assert user claims failed")
	}
	return uc, nil
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

package dto

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonResp struct {
	ginCtx *gin.Context
}

func NewJsonResp(c *gin.Context) *JsonResp {
	return &JsonResp{ginCtx: c}
}

// Success 响应成功
func (j *JsonResp) Success(data ...any) {
	j.responseWithCode(CodeSuccess, data)
}

// FailWithMsg 响应失败，指定消息
func (j *JsonResp) FailWithMsg(code int, msg string) {
	j.response(code, msg, nil)
}

// Fail 响应失败，通过code得到msg
func (j *JsonResp) Fail(code int, err ...error) {
	j.responseWithCode(code, nil)
}

func (j *JsonResp) responseWithCode(code int, data any) {
	var resp = Response{
		Code: code,
		Msg:  getMsgWithCode(code),
		Data: data,
	}
	j.ginCtx.JSON(http.StatusOK, resp)
}

func (j *JsonResp) response(code int, msg string, data any) {
	var resp = Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	j.ginCtx.JSON(http.StatusOK, resp)
}

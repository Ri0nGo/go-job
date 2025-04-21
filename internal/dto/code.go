package dto

import (
	"errors"
	"fmt"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	baseCode   = 10000
	offsetCode = 100
)

const (
	userModule = iota + 1
	jobModule
	nodeModule
	tagModule
	jobRecordModule
)

const (
	CodeSuccess = 0

	// 固定错误码
	UnauthorizedError = 401
	ForbiddenError    = 403
	NotFound          = 404
	ServerError       = 500
)

const (
	// 1000 ~ 1099 公共参数错误码
	ParamsError        = 1000
	EmailFormatError   = 1001
	EmailCodeSendError = 1002

	// 1100 ~ 1199 db 相关错误

	// 1200 ~ 1299 文件相关错误
	FileNotExist      = 1200
	UploadFileError   = 1201
	FileExtNotSupport = 1202
	FileValidError    = 1203
)

var (
	UserNotExist               = genCodeMsg(userModule, 0, "用户不存在")
	UsernameOrPasswordError    = genCodeMsg(userModule, 1, "用户名或密码错误")
	UserAddFailed              = genCodeMsg(userModule, 2, "用户创建失败")
	UserUpdateFailed           = genCodeMsg(userModule, 3, "用户更新失败")
	UserGetFailed              = genCodeMsg(userModule, 4, "用户查询失败")
	UserDeleteFailed           = genCodeMsg(userModule, 5, "用户删除失败")
	UserPasswordComplexRequire = genCodeMsg(userModule, 6, "密码必须包含字母、数字、特殊字符，并且不少于八位")
	UsernameExist              = genCodeMsg(userModule, 7, "用户名已存在")
	UserPasswordNotMatch       = genCodeMsg(userModule, 8, "两次密码不一致")
	UserLoginErr               = genCodeMsg(userModule, 9, "用户登录失败")
)

var (
	JobNotExist     = genCodeMsg(jobModule, 0, "任务不存在")
	JobAddFailed    = genCodeMsg(jobModule, 1, "任务创建失败")
	JobUpdateFailed = genCodeMsg(jobModule, 2, "任务更新失败")
	JobGetFailed    = genCodeMsg(jobModule, 3, "任务查询失败")
	JobDeleteFailed = genCodeMsg(jobModule, 4, "任务删除失败")
)

var (
	JobRecordNotExist     = genCodeMsg(jobRecordModule, 0, "任务记录不存在")
	JobRecordAddFailed    = genCodeMsg(jobRecordModule, 1, "任务记录添加失败")
	JobRecordGetFailed    = genCodeMsg(jobRecordModule, 2, "任务记录查询失败")
	JobRecordDeleteFailed = genCodeMsg(jobRecordModule, 3, "任务记录删除失败")
)

var (
	NodeNotExist     = genCodeMsg(nodeModule, 0, "节点不存在")
	NodeAddFailed    = genCodeMsg(nodeModule, 1, "节点创建失败")
	NodeUpdateFailed = genCodeMsg(nodeModule, 2, "节点更新失败")
	NodeGetFailed    = genCodeMsg(nodeModule, 3, "节点查询失败")
	NodeDeleteFailed = genCodeMsg(nodeModule, 4, "节点删除失败")
)

var msgMap = map[int]string{
	CodeSuccess:       "success",
	ServerError:       "server error",
	UnauthorizedError: "unauthorized",

	// ============= common error code ============= //
	ParamsError:        "params error",
	EmailFormatError:   "邮箱格式错误",
	EmailCodeSendError: "验证码发送失败",

	// ============= file ============= //
	UploadFileError:   "upload file error",
	FileNotExist:      "file not exist",
	FileExtNotSupport: "file ext not support",
	FileValidError:    "file valid error",
}

func getMsgWithCode(code int) string {
	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return ""
}

// genCodeMsg 生成对应的code码
func genCodeMsg(moduleCode, bizCode int, msg string) int {
	c := baseCode + moduleCode*offsetCode + bizCode
	if _, ok := msgMap[c]; ok {
		panic(errors.New(fmt.Sprintf("code: %d 已存在", c)))
	}
	registerCode(c, msg)
	return c
}

func registerCode(code int, msg string) {
	msgMap[code] = msg
}

package dto

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	baseCode   = 1000
	offsetCode = 100
)

const (
	userModule = iota + 1
	jobModule
)

const (
	CodeSuccess = 0

	// 固定错误码
	AuthError         = 401
	UnauthorizedError = 403
	NotFound          = 404
	ServerError       = 500
)

const (
	// 1000 ~ 1009 公共参数错误码
	ParamsError = 1000

	// 1010 ~ 1019 db 相关错误
	// 1030 ~ 1040 文件相关错误
	FileNotExist      = 1030
	UploadFileError   = 1031
	FileExtNotSupport = 1032
	FileValidError    = 1033
)

var (
	UserNotExist = genCodeMsg(userModule, 0, "用户不存在")
)

var (
	JobNotExist     = genCodeMsg(jobModule, 0, "任务不存在")
	JobAddFailed    = genCodeMsg(jobModule, 1, "任务创建失败")
	JobUpdateFailed = genCodeMsg(jobModule, 2, "任务更新失败")
	JobGetFailed    = genCodeMsg(jobModule, 3, "任务查询失败")
	JobDeleteFailed = genCodeMsg(jobModule, 4, "任务删除失败")
)

var msgMap = map[int]string{
	CodeSuccess:       "success",
	AuthError:         "auth error",
	NotFound:          "not found",
	ServerError:       "server error",
	UnauthorizedError: "unauthorized",

	// ============= request ============= //
	ParamsError: "params error",

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
	registerCode(c, msg)
	return c
}

func registerCode(code int, msg string) {
	msgMap[code] = msg
}

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

	AuthError         = 401
	UnauthorizedError = 403
	NotFound          = 404
	ServerError       = 500

	// 1000 ~ 1099 公共错误码
	ParamsError = 1000
)

var (
	UserNotExist = genCodeMsg(userModule, 0, "用户不存在")
)

var (
	JobNotExist  = genCodeMsg(jobModule, 0, "用户不存在")
	JobAddFailed = genCodeMsg(jobModule, 1, "任务创建失败")
)

var msgMap = map[int]string{
	CodeSuccess:       "success",
	AuthError:         "auth error",
	NotFound:          "not found",
	ServerError:       "server error",
	UnauthorizedError: "unauthorized",
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

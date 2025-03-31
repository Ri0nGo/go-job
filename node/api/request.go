package api

type ReqJob struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`           // 任务名称
	ExecType     ExecType `json:"exec_type"`      // 任务类型
	Crontab      string   `json:"crontab"`        // crontab 表达式
	Tags         []string `json:"tags"`           // 标签
	CallBackHttp string   `json:"call_back_http"` // 回调数据的http
}

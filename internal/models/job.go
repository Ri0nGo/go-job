package models

type Job struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`      // 任务名称
	ExecType ExecType `json:"exec_type"` // 任务类型
	Crontab  string   `json:"crontab"`   // crontab 表达式
	Tags     []string `json:"tags"`      // 标签
}

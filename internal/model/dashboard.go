package model

type DashboardKey string

const (
	DashboardKeyDayStatus DashboardKey = "day"
	DashboardKeyJobStatus DashboardKey = "job"
)

type ReqDashboardChart struct {
	Begin int64        `json:"begin" form:"begin"`
	End   int64        `json:"end" form:"end"`
	Key   DashboardKey `json:"key" form:"key"`
}

package dto

type RespDashboardTotalData struct {
	Node NodeDashboard `json:"node"`
	Job  JobDashboard  `json:"job"`
}

type NodeDashboard struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type JobDashboard struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	InActive int `json:"inactive"`
}

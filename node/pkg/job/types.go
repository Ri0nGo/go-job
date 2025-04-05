package job

type JobStatus int

const (
	Pending JobStatus = iota
	Running
	Success
	Failed
)

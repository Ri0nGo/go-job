package paths

import "fmt"

type JobAPI struct {
	BasePath   string
	Create     string
	Update     string
	GetAll     string
	Upload     string
	GetOneById func(id int) string
	DeleteById func(id int) string
}

/*const (
	ApiNodeBase      = "/api/go-job/node"
	ApiNodeCreateJob = ApiNodeBase + "/jobs/add" // POST
	ApiNodeDeleteJob = ApiNodeBase + "/jobs/%d"  // DELETE, fmt.Sprintf(ApiNodeJobByID, jobID)
	ApiNodeUpdateJob = ApiNodeBase + "/jobs"     // PUT
	ApiNodeGetJobs   = ApiNodeBase + "/jobs"     // GET
)*/

var NodeJobAPI = &JobAPI{
	BasePath: "/api/go-job/node/jobs",
	Create:   "/add",
	Update:   "",
	GetAll:   "",
	Upload:   "/upload",
	GetOneById: func(id int) string {
		return fmt.Sprintf("%d", id)
	},
	DeleteById: func(id int) string {
		return fmt.Sprintf("%d", id)
	},
}

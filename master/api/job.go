package api

import (
	"github.com/gin-gonic/gin"
	"go-job/master/service"
)

type JobApi struct {
	JobService service.IJobService
}

func NewJobApi(jobService service.IJobService) *JobApi {
	return &JobApi{
		JobService: jobService,
	}
}

func (jobApi *JobApi) RegisterRoutes(group *gin.RouterGroup) {
	jobGroup := group.Group("/jobs")
	{
		jobGroup.GET(":id")
	}
}

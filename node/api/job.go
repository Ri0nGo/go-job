package api

import (
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/node/service"
	//"github.com/robfig/cron/v3"
)

type JobHandler struct {
	JobService service.IJobService
}

func NewJobHandler(jobService service.IJobService) *JobHandler {
	return &JobHandler{
		JobService: jobService,
	}
}

func (h *JobHandler) RegisterRoutes(server *gin.RouterGroup) {
	jh := server.Group("/job")
	jh.POST("/add", h.AddJob)
}

func (h *JobHandler) AddJob(ctx *gin.Context) {
	var req ReqJob
	if err := ctx.ShouldBindJSON(&req.JobDAO); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.CodeSuccess)
		return
	}

	err := h.JobService.AddJob(ctx.Request.Context(), req.JobDAO)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.JobAddFailed)
		return
	}

	dto.NewJsonResp(ctx).Success()
}

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type JobRecordApi struct {
	JobRecordService service.IJobRecordService
}

func NewJobRecordApi(jobRecordService service.IJobRecordService) *JobRecordApi {
	return &JobRecordApi{
		JobRecordService: jobRecordService,
	}
}

// RegisterRoutes 注册任务记录模块路由
func (a *JobRecordApi) RegisterRoutes(group *gin.RouterGroup) {
	jobRecordGroup := group.Group("/job_records")
	{
		jobRecordGroup.GET("", a.GetJobRecordList)
		jobRecordGroup.GET("/:id", a.GetJobRecord)
		jobRecordGroup.POST("/add", a.AddJobRecord)
		jobRecordGroup.DELETE("/:id", a.DeleteJobRecord)
	}
}

// GetJobRecord 创建任务记录
func (a *JobRecordApi) GetJobRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	jobRecord, err := a.JobRecordService.GetJobRecord(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dto.NewJsonResp(ctx).Success()
		return
	}
	if err != nil {
		slog.Error("get jobRecord err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(jobRecord)
}

// AddJobRecord 添加job记录
func (a *JobRecordApi) AddJobRecord(ctx *gin.Context) {
	var req model.CallbackJobResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("add job record err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	err := a.JobRecordService.AddJobRecord(req)
	if err != nil {
		slog.Error("add job record err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobRecordAddFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// GetJobRecordList 查询任务记录列表
func (a *JobRecordApi) GetJobRecordList(ctx *gin.Context) {
	var req dto.ReqJobRecords
	if err := ctx.ShouldBindQuery(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	list, err := a.JobRecordService.GetJobRecordList(req.Page, req.JobId)
	if err != nil {
		slog.Error("get jobRecord list err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobRecordGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(list)
}

// DeleteJobRecord 删除任务记录
func (a *JobRecordApi) DeleteJobRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	if err := a.JobRecordService.DeleteJobRecord(id); err != nil {
		slog.Error("delete jobRecord err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobRecordDeleteFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

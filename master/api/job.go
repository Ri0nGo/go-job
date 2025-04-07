package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/upload"
	"go-job/master/pkg/config"
	"go-job/master/service"
	"go-job/node/pkg/utils"
	"log/slog"
	"strconv"
	"time"
)

type JobApi struct {
	JobService service.IJobService
}

func NewJobApi(jobService service.IJobService) *JobApi {
	return &JobApi{
		JobService: jobService,
	}
}

func (a *JobApi) RegisterRoutes(group *gin.RouterGroup) {
	jobGroup := group.Group("/jobs")
	{
		jobGroup.GET("", a.GetJob)
		jobGroup.GET("/:id", a.GetJob)
		jobGroup.POST("", a.AddJob)
		jobGroup.PUT("", a.UpdateJob)
		jobGroup.DELETE("/:id", a.DeleteJob)
		jobGroup.POST("/upload", a.UploadFile)
	}
}

func (a *JobApi) GetJob(ctx *gin.Context) {
	var req dto.ReqId
	if err := ctx.ShouldBindQuery(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	job, err := a.JobService.GetJob(req.Id)
	if err != nil {
		slog.Error("get job err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(job)
}

func (a *JobApi) GetJobList(ctx *gin.Context) {
	var page model.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	list, err := a.JobService.GetJobList(page)
	if err != nil {
		slog.Error("get job list err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(list)
}

func (a *JobApi) AddJob(ctx *gin.Context) {
	var req model.Job
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	err := a.JobService.AddJob(req)
	if err != nil {
		slog.Error("add job err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobAddFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

func (a *JobApi) DeleteJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	if err := a.JobService.DeleteJob(id); err != nil {
		slog.Error("delete job err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobDeleteFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

func (a *JobApi) UpdateJob(ctx *gin.Context) {
	var req model.Job
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
	}
	if err := a.JobService.UpdateJob(req); err != nil {
		slog.Error("update job err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobUpdateFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// UploadFile 保存上传的文件(master 用的)
func (a *JobApi) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	fileMeta := upload.FileMeta{
		Filename:     file.Filename,
		UUIDFileName: uuid.New().String(),
		Size:         int(file.Size),
		Uploaded:     time.Now(),
	}
	if err = upload.ValidatorFileOpts(fileMeta,
		upload.FileExtValidator, upload.FileSizeValidator); err != nil {
		dto.NewJsonResp(ctx).FailWithMsg(dto.FileValidError, err.Error())
		return
	}

	if err = utils.EnsureDir(config.App.Data.UploadJobDir); err != nil {
		slog.Error("file dir create error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	// 保存文件
	savePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir, fileMeta.UUIDFileName)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		slog.Error("save file error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	upload.SetFileMeta(fileMeta.UUIDFileName, fileMeta)

	dto.NewJsonResp(ctx).Success(map[string]string{
		"key": fileMeta.UUIDFileName,
	})
}

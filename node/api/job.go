package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/pkg/utils"
	"go-job/internal/upload"
	"go-job/node/pkg/config"
	"go-job/node/service"
	"log/slog"
	"strconv"
	"time"
)

type JobHandler struct {
	JobService service.IJobService
}

func NewJobHandler(jobService service.IJobService) *JobHandler {
	return &JobHandler{
		JobService: jobService,
	}
}

// RegisterRoutes 注册job相关的路由, 遵循restful 风格
func (h *JobHandler) RegisterRoutes(server *gin.RouterGroup) {
	jh := server.Group("/jobs")
	jh.POST("/add", h.AddJob)
	jh.DELETE("/:id", h.DeleteJob)
	jh.PUT("", h.UpdateJob)
	jh.GET("", h.GetJob)
	jh.POST("/upload", h.UploadFile)
	//jh.GET("", h.GetJobList)  todo 待实现
}

// AddJob 添加任务
func (h *JobHandler) AddJob(ctx *gin.Context) {
	var req dto.ReqJob
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("add job bind json err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	err := h.JobService.AddJob(ctx.Request.Context(), req)
	if err != nil {
		slog.Error("add job error", "req", req, "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobAddFailed)
		return
	}

	dto.NewJsonResp(ctx).Success()
}

func (h *JobHandler) DeleteJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	h.JobService.DeleteJob(ctx.Request.Context(), id)
	dto.NewJsonResp(ctx).Success()
}

func (h *JobHandler) UpdateJob(ctx *gin.Context) {
	var req dto.ReqJob
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	err := h.JobService.UpdateJob(ctx.Request.Context(), req)
	if err != nil {
		slog.Error("update job error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobUpdateFailed)
		return
	}

	dto.NewJsonResp(ctx).Success()
}

func (h *JobHandler) GetJob(ctx *gin.Context) {
	var req dto.ReqId
	if err := ctx.ShouldBindQuery(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	job, err := h.JobService.GetJob(ctx.Request.Context(), req.Id)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.JobNotExist)
		return
	}
	dto.NewJsonResp(ctx).Success(dto.RespNodeJob{
		Id:            job.JobMeta.Id,
		Name:          job.JobMeta.Name,
		ExecType:      job.JobMeta.ExecType,
		RunningStatus: job.RunningStatus,
		FileName:      job.JobMeta.FileName,
	})

}

func (h *JobHandler) GetJobList(ctx *gin.Context) {
	dto.NewJsonResp(ctx).Success()
}

// UploadFile 接收文件
// 接收master发送过来的文件，然后保存到当前节点中
func (h *JobHandler) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	uuidFilename, ok := ctx.GetPostForm("filename")
	if !ok {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	// 校验文件信息
	fileMeta := upload.FileMeta{
		Filename:     file.Filename,
		UUIDFileName: uuidFilename,
		Size:         int(file.Size),
		UploadTime:   time.Now(),
	}
	if err := upload.ValidatorFileOpts(fileMeta,
		upload.FileExtValidator, upload.FileSizeValidator); err != nil {
	}

	if err = utils.EnsureDir(config.App.Data.UploadJobDir); err != nil {
		slog.Error("file dir create error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	savePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir, uuidFilename)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		slog.Error("save file error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	dto.NewJsonResp(ctx).Success()
}

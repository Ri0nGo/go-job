package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/upload"
	"go-job/node/pkg/config"
	"go-job/node/pkg/utils"
	"go-job/node/service"
	"log/slog"
	"path/filepath"
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

// RegisterRoutes 注册job相关的路由
func (h *JobHandler) RegisterRoutes(server *gin.RouterGroup) {
	jh := server.Group("/job")
	jh.POST("/add", h.AddJob)
}

// AddJob 添加任务
func (h *JobHandler) AddJob(ctx *gin.Context) {
	var req dto.ReqJob
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	err := h.JobService.AddJob(ctx.Request.Context(), req)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.JobAddFailed)
		return
	}

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
	uuid, ok := ctx.GetPostForm("uuid")
	if !ok {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	// 校验文件信息
	fileMeta := upload.FileMeta{
		Filename:     file.Filename,
		UUIDFileName: uuid,
		Size:         int(file.Size),
		Uploaded:     time.Now(),
	}
	if err := upload.ValidatorFileOpts(fileMeta,
		upload.FileExtValidator, upload.FileSizeValidator); err != nil {
	}

	if err = utils.EnsureDir(config.App.Data.UploadJobDir); err != nil {
		slog.Error("file dir create error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	ext := filepath.Ext(file.Filename)
	savePath := fmt.Sprintf("%s/%s.%s", config.App.Data.UploadJobDir, uuid, ext)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		slog.Error("save file error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	dto.NewJsonResp(ctx).Success()
}

// UploadFile 保存上传的文件(master 用的)
/*func (h *JobHandler) UploadFileBak(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	// todo 这里后续如果新增了其他的错误，需要额外调整响应
	if err = h.validJobUploadFile(file); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.FileExtNotSupport)
		return
	}

	if err = utils.EnsureDir(config.App.Data.TmpData); err != nil {
		slog.Error("file dir create error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	// 保存文件
	savePath := fmt.Sprintf("%s/%s", config.App.Data.TmpData, file.Filename)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		slog.Error("save file error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	tempKey := uuid.New().String()
	upload.SetFileMeta(tempKey, upload.FileMeta{
		Filename: file.Filename,
		Filepath: savePath,
		Uploaded: time.Now(),
	})

	dto.NewJsonResp(ctx).Success(map[string]string{
		"key": tempKey,
	})
}*/

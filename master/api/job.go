package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/internal/upload"
	"go-job/master/pkg/config"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"
)

type JobApi struct {
	JobService service.IJobService
	userSvc    service.IUserService
}

func NewJobApi(jobService service.IJobService, userSvc service.IUserService) *JobApi {
	return &JobApi{
		JobService: jobService,
		userSvc:    userSvc,
	}
}

func (a *JobApi) RegisterRoutes(group *gin.RouterGroup) {
	jobGroup := group.Group("/jobs")
	{
		jobGroup.GET("", a.GetJobList)
		jobGroup.GET("/:id", a.GetJob)
		jobGroup.POST("/add", a.AddJob)
		jobGroup.PUT("/update", a.UpdateJob)
		jobGroup.DELETE("/:id", a.DeleteJob)
		jobGroup.POST("/upload", a.UploadFile)
	}
}

func (a *JobApi) GetJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	job, err := a.JobService.GetJob(uc.Uid, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dto.NewJsonResp(ctx).Success()
		return
	}
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
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	list, err := a.JobService.GetJobList(uc.Uid, page)
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
		slog.Error("add job param err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	req.UserId = uc.Uid

	if err := a.JobService.AddJob(req); err != nil {
		slog.Error("add job error", "err", err)
		if service.IsRespErr(err) {
			dto.NewJsonResp(ctx).FailWithMsg(dto.JobAddFailed, err.Error())
		} else {
			dto.NewJsonResp(ctx).Fail(dto.JobAddFailed)
		}
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
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	if err := a.JobService.DeleteJob(uc.Uid, id); err != nil {
		slog.Error("delete job err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.JobDeleteFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

func (a *JobApi) UpdateJob(ctx *gin.Context) {
	var req model.Job
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("update job paras err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	req.UserId = uc.Uid

	if err := a.JobService.UpdateJob(req); err != nil {
		slog.Error("update job err:", "err", err)
		if service.IsRespErr(err) {
			dto.NewJsonResp(ctx).FailWithMsg(dto.JobUpdateFailed, err.Error())
		} else {
			dto.NewJsonResp(ctx).Fail(dto.JobUpdateFailed)
		}
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// UploadFile 保存上传的文件(master 用的)
func (a *JobApi) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		slog.Error("job upload file err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	uuidKey := uuid.New().String()
	uuidFileName := uuidKey + filepath.Ext(file.Filename)
	fileMeta := upload.FileMeta{
		Filename:     file.Filename,
		UUIDFileName: uuidFileName,
		Size:         int(file.Size),
		UploadTime:   time.Now().Truncate(time.Second),
	}
	if err = upload.ValidatorFileOpts(fileMeta,
		upload.FileExtValidator, upload.FileSizeValidator); err != nil {
		switch {
		case errors.Is(err, upload.ErrFileTooLarge):
			err = service.ErrFileTooLarge
		case errors.Is(err, upload.ErrFileExtNotSupported):
			err = service.ErrJobExtNotSupport
		}
		dto.NewJsonResp(ctx).FailWithMsg(dto.FileValidError, err.Error())
		return
	}

	if err = utils.EnsureDir(config.App.Data.UploadJobDir); err != nil {
		slog.Error("file dir create error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	// 保存文件
	savePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir,
		uuidFileName)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		slog.Error("save file error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UploadFileError)
		return
	}

	upload.SetFileMeta(uuidKey, fileMeta)

	dto.NewJsonResp(ctx).Success(map[string]string{
		"key": uuidKey,
	})
}

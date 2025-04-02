package api

import (
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
)

func AddJob(ctx *gin.Context) {
	var job ReqJob
	if err := ctx.ShouldBindJSON(&job); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.CodeSuccess)
		return
	}
}

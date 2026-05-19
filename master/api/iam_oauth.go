package api

import (
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/master/service"
	"log/slog"
)

type IAMOAuthApi struct {
	oauthSvc service.IIAMOAuthService
}

func NewIAMOAuthApi(oauthSvc service.IIAMOAuthService) *IAMOAuthApi {
	return &IAMOAuthApi{oauthSvc: oauthSvc}
}

func (a *IAMOAuthApi) RegisterRoutes(group *gin.RouterGroup) {
	oauthGroup := group.Group("/oauth")
	{
		oauthGroup.GET("/info", a.GetOAuthInfo)
		oauthGroup.POST("/login", a.Login)
	}
}

func (a *IAMOAuthApi) GetOAuthInfo(ctx *gin.Context) {
	info, err := a.oauthSvc.GetOAuthInfo(ctx.Query("state"))
	if err != nil {
		slog.Error("get oauth info err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	dto.NewJsonResp(ctx).Success(info)
}

func (a *IAMOAuthApi) Login(ctx *gin.Context) {
	var req dto.ReqOAuthLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("oauth login params err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	user, token, err := a.oauthSvc.Login(ctx.Request.Context(), req.Code)
	if err != nil {
		slog.Error("oauth login err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UserLoginErr)
		return
	}
	ctx.Header("Authorization", token)
	dto.NewJsonResp(ctx).Success(dto.RespOAuthLogin{ID: user.Id})
}

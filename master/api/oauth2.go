package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"go-job/internal/dto"
	"go-job/internal/iface/oauth2"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/master/pkg/config"
	"go-job/master/pkg/oauth2/github"
	"go-job/master/pkg/oauth2/qq"
	"go-job/master/service"
	"log/slog"
	"net/http"
)

type OAuth2Api struct {
	oauth2Svc         map[model.AuthType]oauth2.IOAuth2Service
	userSvc           service.IUserService
	authKey           string
	authStateName     string
	authStateExpire   int
	redirectFrontPath string
}

func NewOAuth2Api(userSvc service.IUserService) *OAuth2Api {
	oa := &OAuth2Api{
		oauth2Svc:         make(map[model.AuthType]oauth2.IOAuth2Service),
		userSvc:           userSvc,
		authKey:           "A73fkfiuwbl92smg@iugjChgWth$s89",
		authStateName:     "oauth2-state",
		authStateExpire:   600,
		redirectFrontPath: config.App.OAuth2[model.AuthTypeGithub.String()].RedirectFrontUrl,
	}

	oa.registryOAuth2Svc()
	return oa
}

func (a *OAuth2Api) registryOAuth2Svc() {
	// 注册github
	a.oauth2Svc[model.AuthTypeGithub] = github.NewOAuth2Service(
		config.App.OAuth2[model.AuthTypeGithub.String()].ClientID,
		config.App.OAuth2[model.AuthTypeGithub.String()].ClientSecret,
		config.App.OAuth2[model.AuthTypeGithub.String()].RedirectURL,
	)

	// 注册qq
	a.oauth2Svc[model.AuthTypeQQ] = qq.NewOAuth2Service(
		config.App.OAuth2[model.AuthTypeQQ.String()].ClientID,
		config.App.OAuth2[model.AuthTypeQQ.String()].ClientSecret,
		config.App.OAuth2[model.AuthTypeQQ.String()].RedirectURL,
	)
}

// RegisterRoutes 注册用户模块路由
func (a *OAuth2Api) RegisterRoutes(group *gin.RouterGroup) {
	ouath2Group := group.Group("/oauth2")
	{
		ouath2Group.GET("/github/authurl", a.GithubAuthURL)
		ouath2Group.Any("/github/callback", a.GithubCallback)

		ouath2Group.GET("/qq/authurl", a.QQAuthURL)
		ouath2Group.Any("/qq/callback", a.QQCallback)
	}
}

// ---------------- github ---------------- //

func (a *OAuth2Api) GithubAuthURL(ctx *gin.Context) {
	state := uuid.New()
	authUrl := a.oauth2Svc[model.AuthTypeGithub].GetAuthUrl(ctx, state)
	err := a.setStateCookie(ctx, state)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
	}
	dto.NewJsonResp(ctx).Success(authUrl)
}

func (a *OAuth2Api) GithubCallback(ctx *gin.Context) {
	// 1. 校验state
	//err := a.verifyState(ctx)
	//if err != nil {
	//	slog.Error("verify state failed", "err", err)
	//	dto.NewJsonResp(ctx).FailWithMsg(dto.UnauthorizedError, "非法请求")
	//	return
	//}

	// 2. 通过code获取userinfo
	code := ctx.Query("code")
	authModel, err := a.oauth2Svc[model.AuthTypeGithub].GetAuthIdentity(ctx, code)
	if err != nil {
		slog.Error("get auth identity error", "err", err)
		dto.NewJsonResp(ctx).FailWithMsg(dto.UnauthorizedError, "认证失败")
		return
	}

	// 3. 注册用户
	user, err := a.userSvc.FindOrCreateByAuthIdentity(authModel)
	if err != nil {
		slog.Error("find or create user failed", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	// 4. 生成jwt token
	token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateUserToken(user)
	if err != nil {
		slog.Error("generate jwt token failed", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}

	redirectFullPath := fmt.Sprintf("%s?t=%s&uid=%d", a.redirectFrontPath, token, user.Id)
	slog.Info("redirectFullPath: " + redirectFullPath)
	ctx.Redirect(http.StatusFound, redirectFullPath)
}

func (a *OAuth2Api) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	builder := auth.NewJwtBuilder(a.authKey)
	token, err := builder.GenerateToken(claims)
	if err != nil {
		return err
	}

	ctx.SetCookie(a.authStateName, token,
		a.authStateExpire, "/api/go-job/oauth2/github/callback",
		"", false, false)
	return nil
}

func (a *OAuth2Api) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(a.authStateName)
	if err != nil {
		return err
	}
	var sc StateClaims
	_, err = auth.NewJwtBuilder(a.authKey).ParseToken(&sc, ck)
	if err != nil {
		return err
	}
	if state != sc.State {
		return fmt.Errorf("state different, receive state: %s, jwt state: %s", state, sc.State)
	}
	return nil
}

// ---------------- qq ---------------- //

func (a *OAuth2Api) QQAuthURL(ctx *gin.Context) {
	state := uuid.New()
	authUrl := a.oauth2Svc[model.AuthTypeQQ].GetAuthUrl(ctx, state)
	err := a.setStateCookie(ctx, state)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
	}
	dto.NewJsonResp(ctx).Success(authUrl)
}

func (a *OAuth2Api) QQCallback(ctx *gin.Context) {
	// 1. 校验state
	//err := a.verifyState(ctx)
	//if err != nil {
	//	slog.Error("verify state failed", "err", err)
	//	dto.NewJsonResp(ctx).FailWithMsg(dto.UnauthorizedError, "非法请求")
	//	return
	//}

	// 2. 通过code获取userinfo
	code := ctx.Query("code")
	authModel, err := a.oauth2Svc[model.AuthTypeQQ].GetAuthIdentity(ctx, code)
	if err != nil {
		slog.Error("get auth identity error", "err", err)
		dto.NewJsonResp(ctx).FailWithMsg(dto.UnauthorizedError, "认证失败")
		return
	}

	// 3. 注册用户
	user, err := a.userSvc.FindOrCreateByAuthIdentity(authModel)
	if err != nil {
		slog.Error("find or create qq user failed", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	// 4. 生成jwt token
	token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateUserToken(user)
	if err != nil {
		slog.Error("generate jwt token failed", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}

	redirectFullPath := fmt.Sprintf("%s?t=%s&uid=%d", a.redirectFrontPath, token, user.Id)
	slog.Info("redirectFullPath: " + redirectFullPath)
	ctx.Redirect(http.StatusFound, redirectFullPath)
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

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
	"net/url"
	"strconv"
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
		authStateName:     "oauth2state",
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
	scene, err := a.verifyScene(ctx.Query("scene"))
	if err != nil {
		slog.Error(err.Error())
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	state := uuid.New()
	authUrl := a.oauth2Svc[model.AuthTypeGithub].GetAuthUrl(ctx, state)
	if err = a.setStateCookie(ctx, state, "/api/go-job/oauth2/github/callback"); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	if err = a.userSvc.SaveState(ctx, state, scene, model.AuthTypeGithub); err != nil {
		slog.Error("save state failed", "err", err.Error())
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	dto.NewJsonResp(ctx).Success(authUrl)
}

func (a *OAuth2Api) GithubCallback(ctx *gin.Context) {
	// 验证state是否有效以及是否已经被使用过
	oauth2State, err := a.verifyState(ctx)
	if err != nil {
		slog.Error("verify github state failed", "err", err)
		ctx.Redirect(http.StatusFound, a.redirectFrontPath)
		return
	}

	// 2. 通过code获取userinfo
	code := ctx.Query("code")
	authModel, err := a.oauth2Svc[model.AuthTypeGithub].GetAuthIdentity(ctx, code)
	if err != nil {
		slog.Error("get auth identity error", "err", err)
		ctx.Redirect(http.StatusFound, a.redirectFrontPath)
		return
	}
	// 从账户安全来的
	if oauth2State.Scene == model.Auth2SceneSettingsPage {
		a.callBackFromSettings(ctx, authModel, oauth2State)
	} else { // 从登录页来的
		a.callbackToLogin(ctx, authModel, oauth2State)
	}
}

func (a *OAuth2Api) setStateCookie(ctx *gin.Context, state, path string) error {
	claims := StateClaims{
		State: state,
	}
	builder := auth.NewJwtBuilder(a.authKey)
	token, err := builder.GenerateToken(claims)
	if err != nil {
		return err
	}

	ctx.SetCookie(a.authStateName, token,
		a.authStateExpire, path,
		"", false, false)
	return nil
}

func (a *OAuth2Api) verifyState(ctx *gin.Context) (model.OAuth2State, error) {
	state := ctx.Query("state")               // 获取github回调的state
	token, err := ctx.Cookie(a.authStateName) // 获取cookie中的state
	if err != nil {
		return model.OAuth2State{}, err
	}
	// 校验state是否为篡改
	if err := a.checkStateValid(state, token); err != nil { // 确保state未被篡改
		return model.OAuth2State{}, err
	}
	// 校验state是否已经被使用过了
	oauth2State, err := a.userSvc.VerifyState(ctx, state)
	if err != nil {
		slog.Error("verify state failed", "state", state, "err", err)
		return model.OAuth2State{}, err
	}
	return oauth2State, nil
}

func (a *OAuth2Api) checkStateValid(state string, token string) error {
	var sc StateClaims
	_, err := auth.NewJwtBuilder(a.authKey).ParseToken(&sc, token)
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
	scene, err := a.verifyScene(ctx.Query("scene"))
	if err != nil {
		slog.Error(err.Error())
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	state := uuid.New()
	authUrl := a.oauth2Svc[model.AuthTypeQQ].GetAuthUrl(ctx, state)
	if err = a.setStateCookie(ctx, state, "/api/go-job/oauth2/qq/callback"); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	if err = a.userSvc.SaveState(ctx, state, scene, model.AuthTypeQQ); err != nil {
		slog.Error("save state failed", "err", err.Error())
		dto.NewJsonResp(ctx).Fail(dto.ServerError)
		return
	}
	dto.NewJsonResp(ctx).Success(authUrl)
}

func (a *OAuth2Api) QQCallback(ctx *gin.Context) {
	// 1. 校验state
	oauth2State, err := a.verifyState(ctx)
	if err != nil {
		slog.Error("verify qq state failed", "err", err)
		ctx.Redirect(http.StatusFound, a.redirectFrontPath)
		return
	}

	// 2. 通过code获取userinfo
	code := ctx.Query("code")
	authModel, err := a.oauth2Svc[model.AuthTypeQQ].GetAuthIdentity(ctx, code)
	if err != nil {
		slog.Error("get auth identity error", "err", err)
		dto.NewJsonResp(ctx).FailWithMsg(dto.UnauthorizedError, "认证失败")
		return
	}

	// 从账户安全来的
	if oauth2State.Scene == model.Auth2SceneSettingsPage {
		a.callBackFromSettings(ctx, authModel, oauth2State)
	} else { // 从登录页来的
		a.callbackToLogin(ctx, authModel, oauth2State)
	}
}

func (a *OAuth2Api) verifyScene(scene string) (model.Auth2Scene, error) {
	validScene := map[string]model.Auth2Scene{
		"account-security": model.Auth2SceneSettingsPage,
		"login":            model.Auth2SceneLoginPage,
	}
	if val, ok := validScene[scene]; ok {
		return val, nil
	}
	return "", fmt.Errorf("scene %s is not valid", scene)
}

func (a *OAuth2Api) callBackFromSettings(ctx *gin.Context, authModel model.AuthIdentity, oauth2State model.OAuth2State) {
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim error", "err", err)
		ctx.Redirect(http.StatusFound, a.redirectFrontPath)
	}
	authModel.UserID = uc.Uid
	if err := a.userSvc.BindOAuth2FromSettings(ctx, authModel); err != nil {
		slog.Error("bind oauth2 error", "err", err)
		redirectPath := fmt.Sprintf("%s?status=%d", oauth2State.RedirectPage, 1)
		ctx.Redirect(http.StatusFound, redirectPath) // 跳转会原来的页面
		return
	}
	redirectPath := fmt.Sprintf("%s?status=%d", oauth2State.RedirectPage, 0)
	ctx.Redirect(http.StatusFound, redirectPath)
}

func (a *OAuth2Api) callbackToLogin(ctx *gin.Context, authModel model.AuthIdentity,
	oauth2State model.OAuth2State) {
	key := uuid.New()
	domainUser, ret, err := a.userSvc.UserOAuth2FromLogin(ctx, key, authModel)
	if err != nil {
		slog.Error("get user OAuth2 error", "err", err)
		ctx.Redirect(http.StatusFound, a.redirectFrontPath)
		return
	}

	q := url.Values{}
	q.Set("redirect_page", oauth2State.RedirectPage)
	switch ret {
	case 0: // 第三方账号未绑定用户
		q.Set("key", key) // 用来获取用户信息
		q.Set("platform", oauth2State.Platform)
		redirectURL := a.redirectFrontPath + "?" + q.Encode()
		slog.Info("bind oauth2 account", "redirectURL", redirectURL)
		ctx.Redirect(302, redirectURL)
	case 1: // 第三方账号已经绑定过用户
		token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateToken(UserDomainToClaim(domainUser))
		if err != nil {
			slog.Error("generate token failed", "err", err)
			ctx.Redirect(http.StatusFound, a.redirectFrontPath)
		}
		q.Set("t", token)
		q.Set("uid", strconv.Itoa(domainUser.Id))
		redirectURL := a.redirectFrontPath + "?" + q.Encode()
		slog.Info("already bind oauth2 account", "redirectURL", redirectURL)
		ctx.Redirect(302, redirectURL)
	default:
		ctx.Redirect(302, a.redirectFrontPath)
	}
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

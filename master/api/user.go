package api

import (
	"context"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/internal/pkg/consts"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/config"
	"go-job/master/repo/cache"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

const (
	// 官方自带的regexp 不能识别 ?= 这样的正则表达式
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	biz                  = "user"
)

type UserApi struct {
	passwordRegex *regexp.Regexp
	userService   service.IUserService
	codeSvc       service.IEmailCodeService
}

func NewUserApi(userService service.IUserService, codeSvc service.IEmailCodeService) *UserApi {
	return &UserApi{
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		userService:   userService,
		codeSvc:       codeSvc,
	}
}

// RegisterRoutes 注册用户模块路由
func (a *UserApi) RegisterRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/users")
	{
		userGroup.GET("", a.GetUserList)
		userGroup.GET("/:id", a.GetUser)
		userGroup.POST("/add", a.AddUser)
		userGroup.PUT("/update", a.UpdateUser)
		userGroup.DELETE("/:id", a.DeleteUser)
		userGroup.POST("/login", a.Login)

		userGroup.POST("/bind/email/code_send", a.BindEmailCodeSend)
		userGroup.POST("/bind/email", a.BindEmail)
	}
}

// GetUser 查询用户
func (a *UserApi) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	user, err := a.userService.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dto.NewJsonResp(ctx).Success()
		return
	}
	if err != nil {
		slog.Error("get user err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UserNotExist)
		return
	}
	dto.NewJsonResp(ctx).Success(user)
}

// GetUserList 查询用户列表
func (a *UserApi) GetUserList(ctx *gin.Context) {
	var page model.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	list, err := a.userService.GetUserList(page)
	if err != nil {
		slog.Error("get user list err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UserGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(list)
}

// AddUser 添加用户
func (a *UserApi) AddUser(ctx *gin.Context) {
	var req model.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	if req.Password != req.ConfirmPassword {
		dto.NewJsonResp(ctx).Fail(dto.UserPasswordNotMatch)
		return
	}

	if !a.isComplexPassword(req.Password) {
		dto.NewJsonResp(ctx).Fail(dto.UserPasswordComplexRequire)
		return
	}
	user, err := a.userService.AddUser(req)
	switch err {
	case nil:
		dto.NewJsonResp(ctx).Success(map[string]int{"id": user.Id})
		return
	case service.ErrUsernameDuplicate:
		dto.NewJsonResp(ctx).Fail(dto.UsernameExist)
		return
	default:
		slog.Error("add user err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UserAddFailed)
		return
	}
}

func (a *UserApi) isComplexPassword(password string) bool {
	ok, err := a.passwordRegex.MatchString(password)
	if err != nil {
		return false
	}
	return ok
}

// DeleteUser 删除用户
func (a *UserApi) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	if err := a.userService.DeleteUser(id); err != nil {
		slog.Error("delete user err:", "err", err)
		// TODO 这里需要区分是查询不到ID还是删除异常了，如果是查询不到ID则直接返回success
		dto.NewJsonResp(ctx).Fail(dto.UserDeleteFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// UpdateUser 更新用户(不含密码)
func (a *UserApi) UpdateUser(ctx *gin.Context) {
	var req model.DomainUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	// 若用户不存在也直接提示更新成功
	/*	if _, err := a.userService.GetUser(req.Id); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.UserNotExist)
		return
	}*/

	if err := a.userService.UpdateUser(req); err != nil {
		slog.Error("update user error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UserUpdateFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// Login 用户登录
func (a *UserApi) Login(ctx *gin.Context) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("login error", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	domainUser, err := a.userService.Login(req.Username, req.Password)
	switch err {
	case nil:
		token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateToken(a.userDomainToClaim(domainUser))
		if err != nil {
			slog.Error("create token err", "err", err)
			dto.NewJsonResp(ctx).Fail(dto.ServerError)
			return
		}
		ctx.Header("Authorization", token)
		dto.NewJsonResp(ctx).Success(map[string]int{
			"id": domainUser.Id,
		})
		return
	case service.ErrInvalidUserOrPassword:
		dto.NewJsonResp(ctx).Fail(dto.UsernameOrPasswordError)
		return
	default:
		dto.NewJsonResp(ctx).Fail(dto.UserLoginErr)
	}
}

// Bind 绑定邮箱
func (a *UserApi) BindEmail(ctx *gin.Context) {
	var req dto.ReqUserEmailBind
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	// 验证邮箱
	if ok := utils.IsValidEmail(req.Email); !ok {
		dto.NewJsonResp(ctx).Fail(dto.EmailFormatError)
		return
	}
	// 验证用户是否合法
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	if _, err = a.userService.GetUser(uc.Uid); err != nil {
		slog.Error("get user err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	if err := a.codeSvc.Verify(ctx, biz, req.Email, req.Code); err != nil {
		slog.Error("verify email bind err", "err", err)
		switch err {
		case cache.ErrInputCodeInvalid:
			dto.NewJsonResp(ctx).FailWithMsg(dto.EmailCodeVerifyError, err.Error())
		case cache.ErrCodeVerifyTooMany:
			dto.NewJsonResp(ctx).FailWithMsg(dto.EmailCodeVerifyError, err.Error())
		default:
			dto.NewJsonResp(ctx).Fail(dto.EmailCodeVerifyError)
		}
		return
	}
	if err = a.userService.UserBind(uc.Uid, req.Email); err != nil {
		slog.Error("bind user email err", "err", err)
		if errors.Is(err, service.ErrDuplicateEmail) {
			dto.NewJsonResp(ctx).FailWithMsg(dto.UserEmailBindErr, err.Error())
			return
		}
		dto.NewJsonResp(ctx).Fail(dto.UserEmailBindErr)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// CodeSend 验证码发送
func (a *UserApi) BindEmailCodeSend(ctx *gin.Context) {
	var req dto.ReqUserEmailBind
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	// 验证邮箱
	if ok := utils.IsValidEmail(req.Email); !ok {
		dto.NewJsonResp(ctx).Fail(dto.EmailFormatError)
		return
	}
	// 验证用户是否合法
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	// TODO 再次通过数据库查询用户是否有必要？
	if _, err = a.userService.GetUser(uc.Uid); err != nil {
		slog.Error("get user err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	if err := a.codeSvc.Send(context.Background(), biz, req.Email); err != nil {
		slog.Error("code send err", "err", err)
		switch {
		case errors.Is(err, cache.ErrCodeSendTooMany):
			dto.NewJsonResp(ctx).FailWithMsg(dto.EmailCodeSendError, err.Error())
		default:
			dto.NewJsonResp(ctx).Fail(dto.EmailCodeSendError)
		}
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// GetUserClaim 获取用户的UC信息
func GetUserClaim(ctx *gin.Context) (*model.UserClaims, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("user not exists in context")
	}
	uc, ok := value.(*model.UserClaims)
	if !ok {
		return nil, errors.New("assert user claims failed")
	}
	return uc, nil
}

func (a *UserApi) userDomainToClaim(user model.DomainUser) model.UserClaims {
	uc := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.DefaultLoginJwtExpireTime)),
		},
		Uid: user.Id,
	}
	return uc
}

package api

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/config"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

const (
	// 官方自带的regexp 不能识别 ?= 这样的正则表达式
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserApi struct {
	passwordRegex *regexp.Regexp
	UserService   service.IUserService
}

func NewUserApi(userService service.IUserService) *UserApi {
	return &UserApi{
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		UserService:   userService,
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

		userGroup.POST("/bind/code/send", a.CodeSend)
		userGroup.POST("/bind", a.Bind)
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
	user, err := a.UserService.GetUser(id)
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
	list, err := a.UserService.GetUserList(page)
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
	user, err := a.UserService.AddUser(req)
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
	if err := a.UserService.DeleteUser(id); err != nil {
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
	/*	if _, err := a.UserService.GetUser(req.Id); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.UserNotExist)
		return
	}*/

	if err := a.UserService.UpdateUser(req); err != nil {
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

	domainUser, err := a.UserService.Login(req.Username, req.Password)
	switch err {
	case nil:
		token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateToken(domainUser)
		if err != nil {
			slog.Error("create token err", "err", err)
			dto.NewJsonResp(ctx).Fail(dto.ServerError)
			return
		}
		ctx.Header("Authorization", token)
		dto.NewJsonResp(ctx).Success()
		return
	case service.ErrInvalidUserOrPassword:
		dto.NewJsonResp(ctx).Fail(dto.UsernameOrPasswordError)
		return
	default:
		dto.NewJsonResp(ctx).Fail(dto.UserLoginErr)
	}
}

// Bind 绑定邮箱
func (a *UserApi) Bind(ctx *gin.Context) {
	var req dto.ReqUserBind
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

}

// CodeSend 验证码发送
func (a *UserApi) CodeSend(ctx *gin.Context) {
	var req dto.ReqCodeSend
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
	uc, err := a.getUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}
	if _, err = a.UserService.GetUser(uc.Uid); err != nil {
		slog.Error("get user err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	if err := a.UserService.UserCodeSend(req); err != nil {
		slog.Error("code send err", "err", err)
	}

}

// getUserClaim 获取用户的UC信息
func (a *UserApi) getUserClaim(ctx *gin.Context) (*model.UserClaims, error) {
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

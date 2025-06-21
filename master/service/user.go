package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go-job/internal/dto"
	"go-job/internal/iface/oauth2"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/config"
	"go-job/master/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrUsernameDuplicate     = errors.New("duplicate username")
	ErrInvalidUserOrPassword = errors.New("用户名或密码不对")
	ErrDuplicateEmail        = errors.New("邮箱已被占用")
	ErrOAuth2KeyIsExpired    = errors.New("授权时间过长，请重新授权")
	ErrOAuth2UserDuplicate   = errors.New("该账号已经绑定过")
)

var (
	accountSecurityPage = "/settings/profile"
	userBindOAuth2Page  = "/oauth2/bind"
)

type IUserService interface {
	GetUser(id int) (model.DomainUser, error)
	Login(username string, password string) (model.DomainUser, error)
	GetUserList(page model.Page) (model.Page, error)
	AddUser(user model.User) (model.User, error)
	DeleteUser(id int) error
	UpdateUser(user model.DomainUser) error
	UserSecurity(uid int) (dto.RespUserSecurity, error)

	UserBind(id int, email string) error
	SaveState(ctx context.Context, uid int, state string, scene model.Auth2Scene, platform model.AuthType) error
	VerifyState(ctx context.Context, state string) (model.OAuth2State, error)

	SaveOAuth2Code(ctx context.Context, code string, tempCode model.OAuth2TempCode)
	OAuth2Bind(ctx context.Context, req dto.ReqOAuth2Bind) (model.DomainUser, error)
	OAuth2UnBind(ctx context.Context, uid int, req dto.ReqOAuth2UnBind) error

	OAuth2Code(ctx context.Context, code string) dto.RespOAuth2Code
}

type UserService struct {
	UserRepo    repo.IUserRepo
	oauth2Cache oauth2.IOAuth2Cache
}

func (s *UserService) GetUser(id int) (model.DomainUser, error) {
	if id <= 0 {
		return model.DomainUser{}, errors.New("user id is zero")
	}
	userDao, err := s.UserRepo.QueryById(id)
	if err != nil {
		return model.DomainUser{}, err
	}
	return s.userToDomainUser(userDao), nil
}

func (s *UserService) GetUserList(page model.Page) (model.Page, error) {
	return s.UserRepo.QueryList(page)
}

func (s *UserService) AddUser(user model.User) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hash)

	err = s.UserRepo.Insert(&user)
	var me *mysql.MySQLError
	switch {
	case err == nil:
		return user, nil
	case errors.As(err, &me):
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return user, ErrUsernameDuplicate
		}
		return user, err
	default:
		return user, err
	}
}

func (s *UserService) DeleteUser(id int) error {
	return s.UserRepo.Delete(id)
}

func (s *UserService) UpdateUser(user model.DomainUser) error {
	return s.UserRepo.Update(s.domainUserToUser(user))
}

// SaveState 将state信息存储到redis 中
func (s *UserService) SaveState(ctx context.Context, uid int, state string,
	scene model.Auth2Scene, platform model.AuthType) error {
	return s.oauth2Cache.SetState(ctx, state, model.OAuth2State{
		Uid:      uid,
		State:    state,
		Scene:    scene,
		Platform: platform.String(),
		Used:     false,
	})
}

func (s *UserService) VerifyState(ctx context.Context, state string) (model.OAuth2State, error) {
	auth2State, err := s.oauth2Cache.GetState(ctx, state)
	if err != nil { // 获取state 失败
		return model.OAuth2State{}, err
	}
	if auth2State.Used == true {
		return auth2State, errors.New("state is already used")
	}
	return auth2State, nil
}

func (s *UserService) Login(username, password string) (model.DomainUser, error) {
	user, err := s.UserRepo.QueryByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.userToDomainUser(user), ErrInvalidUserOrPassword
	}
	if err != nil {
		return s.userToDomainUser(user), err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return s.userToDomainUser(user), ErrInvalidUserOrPassword
	}

	if err = s.UserRepo.UpdateDataById(user.Id, map[string]any{
		"login_time": time.Now(),
	}); err != nil {
		return s.userToDomainUser(user), err
	}
	return s.userToDomainUser(user), nil
}

func (s *UserService) UserBind(id int, email string) error {
	err := s.UserRepo.Update(model.User{Id: id, Email: &email})
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}
	return err
}

// OAuth2Bind 绑定第三方账户
func (s *UserService) OAuth2Bind(ctx context.Context, req dto.ReqOAuth2Bind) (model.DomainUser, error) {
	// 查询code 是否过期
	val, err := s.oauth2Cache.GetAuth(ctx, req.Code)
	if err != nil {
		slog.Error("get auth key failed", "key", req.Code, "err", err)
		return model.DomainUser{}, ErrOAuth2KeyIsExpired
	}

	codeData := MapToTempCode(val)
	if codeData.Err != "" {
		return model.DomainUser{}, errors.New(codeData.Err)
	}
	authType, err := platformToAuthType(codeData.Platform)
	if err != nil {
		slog.Error("auth type covert to int failed", "platform", codeData.Platform, "err", err)
		return model.DomainUser{}, err
	}

	// 验证用户名和密码
	domainUser, err := s.Login(req.Username, req.Password)
	if err != nil {
		return model.DomainUser{}, err
	}

	// 校验用户是否已经绑定过
	authModel, err := s.UserRepo.QueryAuth(authType, codeData.Identify)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // 查询异常了
		return model.DomainUser{}, err
	}
	if authModel.ID != 0 { // 该账号已经绑定过了
		return model.DomainUser{}, fmt.Errorf("%w%s", ErrOAuth2UserDuplicate, authType.String())
	}

	if err = s.UserRepo.CreateAuth(&model.AuthIdentity{
		UserID:   domainUser.Id,
		Type:     authType,
		Identity: codeData.Identify,
		Name:     codeData.Name,
	}); err != nil {
		return model.DomainUser{}, err
	}
	return domainUser, nil
}

func (s *UserService) OAuth2UnBind(ctx context.Context, uid int, req dto.ReqOAuth2UnBind) error {
	return s.UserRepo.DeleteAuthByUid(uid, req.AuthType)
}

func (s *UserService) SaveOAuth2Code(ctx context.Context, code string, tempCode model.OAuth2TempCode) {
	if err := s.oauth2Cache.SetAuth(ctx, code, tempCodeToMap(tempCode), time.Minute*5); err != nil {
		slog.Error("save oauth2 code failed", "err", err)
	}
}

func (s *UserService) OAuth2Code(ctx context.Context, code string) dto.RespOAuth2Code {
	var result dto.RespOAuth2Code
	val, err := s.oauth2Cache.GetAuth(ctx, code)
	if err != nil {
		slog.Error("get oauth2 code failed", "err", err)
		return dto.RespOAuth2Code{
			RedirectPage: "/",
			Err:          "认证超时",
		}
	}

	codeData := MapToTempCode(val)
	if codeData.Err != "" {
		result = dto.RespOAuth2Code{
			RedirectPage: "/",
			Err:          codeData.Err,
		}
		return result
	}
	authType, err := platformToAuthType(codeData.Platform)
	if err != nil {
		result = dto.RespOAuth2Code{
			RedirectPage: "/",
			Err:          codeData.Err,
		}
		return result
	}

	switch codeData.Scene {
	case model.Auth2SceneAccountSecurityPage:
		if err = s.handleOAuth2AccountSecurity(authType, codeData); err != nil {
			return dto.RespOAuth2Code{
				RedirectPage: accountSecurityPage,
				Err:          err.Error(),
			}
		}
		result = dto.RespOAuth2Code{
			RedirectPage: accountSecurityPage,
		}
	case model.Auth2SceneLoginPage:
		result = s.handleOAuth2Login(authType, codeData)
	default: // 理论上不可能走到这里
		slog.Error("unknown auth scene", "scene", codeData.Scene)
		result = dto.RespOAuth2Code{
			RedirectPage: "/",
			Err:          codeData.Err,
		}
	}

	result.Platform = codeData.Platform
	return result
}

func tempCodeToMap(tempCode model.OAuth2TempCode) map[string]string {
	return map[string]string{
		"uid":      strconv.Itoa(tempCode.Uid),
		"name":     tempCode.Name,
		"identify": tempCode.Identify,
		"platform": tempCode.Platform,
		"scene":    string(tempCode.Scene),
		"err":      tempCode.Err,
	}
}

func MapToTempCode(val map[string]string) model.OAuth2TempCode {
	uid, err := strconv.Atoi(val["uid"])
	if err != nil {
		slog.Error("convert uid to int failed", "uid", val["uid"])
	}
	return model.OAuth2TempCode{
		Uid:      uid,
		Name:     val["name"],
		Identify: val["identify"],
		Platform: val["platform"],
		Scene:    model.Auth2Scene(val["scene"]),
		Err:      val["err"],
	}
}

func platformToAuthType(platform string) (model.AuthType, error) {
	switch platform {
	case model.AuthTypeQQ.String():
		return model.AuthTypeQQ, nil
	case model.AuthTypeGithub.String():
		return model.AuthTypeGithub, nil
	default:
		return 0, errors.New("unknown auth type")
	}
}

func (s *UserService) handleOAuth2AccountSecurity(authType model.AuthType, codeData model.OAuth2TempCode) error {
	_, err := s.UserRepo.QueryAuth(authType, codeData.Identify)
	switch err {
	case nil: // 第三方账号已经被使用了
		return errors.New("该第三方账号已被使用，无法重复绑定")
	case gorm.ErrRecordNotFound: // 第三方账号没有被使用过
		if err = s.UserRepo.CreateAuth(&model.AuthIdentity{
			UserID:   codeData.Uid,
			Type:     authType,
			Identity: codeData.Identify,
			Name:     codeData.Name,
		}); err != nil {
			slog.Error("save auth identity failed", "err", err)
			return errors.New("系统错误")
		}
		return nil
	default: // 内部错误或其他
		slog.Error("query auth failed", "err", err)
		return errors.New("系统错误")
	}
}

func (s *UserService) handleOAuth2Login(authType model.AuthType, codeData model.OAuth2TempCode) dto.RespOAuth2Code {
	authIdentity, err := s.UserRepo.QueryAuth(authType, codeData.Identify)
	switch err {
	case gorm.ErrRecordNotFound: // 第三方账号未绑定过系统用户
		return dto.RespOAuth2Code{
			RedirectPage: userBindOAuth2Page,
		}
	case nil: // 已绑定过系统用户，直接登录
		token, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateUserToken(model.DomainUser{Id: authIdentity.UserID})
		if err != nil {
			slog.Error("generate user token failed", "err", err)
			return dto.RespOAuth2Code{
				RedirectPage: "/",
				Err:          "系统错误",
			}
		}
		return dto.RespOAuth2Code{
			RedirectPage: "/",
			Token:        token,
		}
	default:
		return dto.RespOAuth2Code{
			RedirectPage: "/",
			Err:          "系统错误",
		}
	}
}

func (s *UserService) UserSecurity(uid int) (dto.RespUserSecurity, error) {
	userAuthInfo, err := s.UserRepo.QueryUserSecurity(uid)
	if err != nil {
		return dto.RespUserSecurity{}, err
	}
	var resp dto.RespUserSecurity
	resp.Email = userAuthInfo.Email
	for _, auth := range userAuthInfo.Auths {
		switch auth.Type {
		case model.AuthTypeQQ:
			resp.QQ = true
		case model.AuthTypeGithub:
			resp.Github = true
		}
	}
	return resp, nil
}

func (s *UserService) domainUserToUser(user model.DomainUser) model.User {
	return model.User{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		About:    user.About,
	}
}

func (s *UserService) userToDomainUser(user model.User) model.DomainUser {
	return model.DomainUser{
		Id:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       utils.PtrToVal(user.Email),
		CreatedTime: user.CreatedTime,
		UpdatedTime: user.UpdatedTime,
		LoginTime:   user.LoginTime,
		About:       user.About,
	}
}

func (s *UserService) getUsername() (string, error) {
	for i := 0; i < 5; i++ {
		userName := s.generateNotDuplicatedUserName()
		_, err := s.UserRepo.QueryByUsername(userName)
		switch err {
		case gorm.ErrRecordNotFound: // 用户名可以用
			return userName, nil
		case nil: // 用户名称重复了
		default: // 数据库出问题了
			return "", err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return "", errors.New("生成用户名称失败")
}

func (s *UserService) generateNotDuplicatedUserName() string {
	min := int64(10000000000) // 11位
	max := int64(99999999999)
	randomId := rand.Int63n(max-min+1) + min
	return "用户" + strconv.Itoa(int(randomId))
}

func NewUserService(userRepo repo.IUserRepo, oauth2Cache oauth2.IOAuth2Cache) IUserService {
	return &UserService{
		UserRepo:    userRepo,
		oauth2Cache: oauth2Cache,
	}
}

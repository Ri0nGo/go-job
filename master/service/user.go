package service

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/rogpeppe/go-internal/cache"
	"go-job/internal/iface/oauth2"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrUsernameDuplicate     = errors.New("duplicate username")
	ErrInvalidUserOrPassword = errors.New("用户名或密码不对")
	ErrDuplicateEmail        = errors.New("邮箱已被占用")
)

var sceneRedirectPath = map[model.Auth2Scene]string{
	model.Auth2SceneSettingsPage: "/account/settings/security",
	model.Auth2SceneLoginPage:    "/",
}

type IUserService interface {
	GetUser(id int) (model.DomainUser, error)
	Login(username string, password string) (model.DomainUser, error)
	GetUserList(page model.Page) (model.Page, error)
	AddUser(user model.User) (model.User, error)
	DeleteUser(id int) error
	UpdateUser(user model.DomainUser) error

	UserBind(id int, email string) error
	SaveState(ctx context.Context, state string, scene model.Auth2Scene, platform model.AuthType) error
	VerifyState(ctx context.Context, state string) (model.OAuth2State, error)

	// oauth2
	FindOrCreateByAuthIdentity(identity model.AuthIdentity) (model.DomainUser, error)

	// UserOAuth2FromLogin 入参，key：认证唯一标识
	// 返回 用户实例，状态(0：用户不存在，1：存在，2 位置错误)，具体错误
	UserOAuth2FromLogin(ctx context.Context, key string, identity model.AuthIdentity) (model.DomainUser, int, error)
	// BindOAuth2 绑定OAuth2
	BindOAuth2FromSettings(ctx context.Context, authModel model.AuthIdentity) error
}

type UserService struct {
	UserRepo repo.IUserRepo
	redisSvc cache.Cache
	stateSvc oauth2.IOAuth2StateCache
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

func (s *UserService) userToDomainUser(user model.User) model.DomainUser {
	return model.DomainUser{
		Id:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       utils.PtrToVal(user.Email),
		CreatedTime: user.CreatedTime,
		About:       user.About,
	}
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
func (s *UserService) SaveState(ctx context.Context, state string,
	scene model.Auth2Scene, platform model.AuthType) error {
	return s.stateSvc.SetState(ctx, state, model.OAuth2State{
		State:        state,
		Scene:        scene,
		RedirectPage: sceneRedirectPath[scene],
		Platform:     platform.String(),
		Used:         false,
	})
}

func (s *UserService) VerifyState(ctx context.Context, state string) (model.OAuth2State, error) {
	auth2State, err := s.stateSvc.GetState(ctx, state)
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

func (s *UserService) FindOrCreateByAuthIdentity(authModel model.AuthIdentity) (model.DomainUser, error) {
	// 1. 查询之前是否登录过
	user, err := s.UserRepo.QueryByAuth(authModel.Type, authModel.Identity)
	switch err {
	case gorm.ErrRecordNotFound:
		// 创建用户
		userName, err := s.getUsername()
		if err != nil {
			return model.DomainUser{}, err
		}
		userModel := &model.User{
			Username: userName,
			Nickname: authModel.Name,
		}
		if err = s.UserRepo.CreateUserAndAuth(userModel, &authModel); err != nil {
			return model.DomainUser{}, err
		}
		return s.userToDomainUser(*userModel), nil
	case nil:
		return s.userToDomainUser(user), nil
	default:
		return model.DomainUser{}, err
	}
}

// UserOAuthByAuthIdentity
func (s *UserService) UserOAuth2FromLogin(ctx context.Context, key string, identity model.AuthIdentity) (model.DomainUser, int, error) {
	user, err := s.UserRepo.QueryByAuth(identity.Type, identity.Identity)
	switch err {
	case gorm.ErrRecordNotFound: // 第三方账户未绑定已存在用户
		value := map[string]string{
			"identity": identity.Identity,
			"type":     strconv.Itoa(int(identity.Type)),
			"name":     identity.Name,
		}
		// 将获取的用户信息存储到redis，方便用户绑定
		if err = s.stateSvc.SetAuth(ctx, key, value, time.Minute*5); err != nil {
			return model.DomainUser{}, 0, err
		}
		return model.DomainUser{}, 0, nil
	case nil: // 用户已存在，则直接登录
		return s.userToDomainUser(user), 1, nil
	default:
		return model.DomainUser{}, 2, err
	}
}

func (s *UserService) BindOAuth2FromSettings(ctx context.Context, authModel model.AuthIdentity) error {
	return s.UserRepo.CreateAuth(&authModel)
}

func (s *UserService) domainUserToUser(user model.DomainUser) model.User {
	return model.User{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		About:    user.About,
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

func NewUserService(userRepo repo.IUserRepo) IUserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

package service

import (
	"context"
	"errors"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/master/pkg/config"
	"go-job/master/repo"
	"log/slog"
	"strings"
	"time"

	oauth2Client "github.com/Ri0nGo/gokit/iam/oauth2Client"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const oauthOpenIDKeyPrefix = "go-job:oauth:openid:"

var (
	ErrOAuthDisabled           = errors.New("OAuth2未启用")
	ErrOAuthClientNotReady     = errors.New("OAuth2客户端未初始化")
	ErrOAuthTokenMissingOpenID = errors.New("OAuth2 token缺少access_token或openid")
)

type IIAMOAuthService interface {
	GetOAuthInfo(state string) (dto.OAuthInfo, error)
	Login(ctx context.Context, code string) (model.DomainUser, string, error)
}

type IAMOAuthService struct {
	enabled     bool
	clientID    string
	redirectURI string
	scopes      []string
	redis       redis.Cmdable
	userRepo    repo.IUserRepo
	client      *oauth2Client.Client
}

func NewIAMOAuthService(redis redis.Cmdable, userRepo repo.IUserRepo) IIAMOAuthService {
	oauthCfg := config.App.OAuth
	scopes := getIAMOAuthScopes(oauthCfg)

	var client *oauth2Client.Client
	if oauthCfg.Enabled {
		var err error
		client, err = oauth2Client.New(oauth2Client.Config{
			BaseURL:      oauthCfg.AuthBaseURL,
			ClientID:     oauthCfg.ClientID,
			ClientSecret: oauthCfg.ClientSecret,
			RedirectURI:  oauthCfg.RedirectURI,
			Scopes:       scopes,
			Timeout:      time.Duration(oauthCfg.TimeoutSeconds) * time.Second,
		})
		if err != nil {
			panic(err)
		}
	}

	return &IAMOAuthService{
		enabled:     oauthCfg.Enabled,
		clientID:    oauthCfg.ClientID,
		redirectURI: oauthCfg.RedirectURI,
		scopes:      scopes,
		redis:       redis,
		userRepo:    userRepo,
		client:      client,
	}
}

func (s *IAMOAuthService) GetOAuthInfo(state string) (dto.OAuthInfo, error) {
	info := dto.OAuthInfo{
		Enabled:      s.enabled,
		ClientID:     s.clientID,
		RedirectURI:  s.redirectURI,
		ResponseType: "code",
		Scope:        strings.Join(s.scopes, " "),
	}
	if !s.enabled {
		return info, nil
	}
	if s.client == nil {
		return dto.OAuthInfo{}, ErrOAuthClientNotReady
	}

	authURL, err := s.client.AuthCodeURL(state)
	if err != nil {
		return dto.OAuthInfo{}, err
	}
	info.AuthURL = authURL
	return info, nil
}

func (s *IAMOAuthService) Login(ctx context.Context, code string) (model.DomainUser, string, error) {
	if err := s.ensureEnabled(); err != nil {
		return model.DomainUser{}, "", err
	}
	token, err := s.client.GetUserToken(ctx, code)
	if err != nil {
		return model.DomainUser{}, "", err
	}
	if err = s.saveOpenID(ctx, token); err != nil {
		return model.DomainUser{}, "", err
	}

	userInfo, err := s.client.GetUserInfo(ctx, token.AccessToken, token.OpenID)
	if err != nil {
		return model.DomainUser{}, "", err
	}
	domainUser, err := s.getOrCreateUserByIAM(userInfo)
	if err != nil {
		return model.DomainUser{}, "", err
	}
	if err = s.userRepo.UpdateLoginTimeByid(domainUser.Id); err != nil {
		return model.DomainUser{}, "", err
	}
	domainUser.LoginTime = time.Now()

	jwtToken, err := auth.NewJwtBuilder(config.App.Server.Key).GenerateUserToken(domainUser)
	if err != nil {
		return model.DomainUser{}, "", err
	}
	return domainUser, jwtToken, nil
}

func (s *IAMOAuthService) getOrCreateUserByIAM(userInfo *oauth2Client.UserInfo) (model.DomainUser, error) {
	if userInfo == nil || strings.TrimSpace(userInfo.OpenID) == "" {
		return model.DomainUser{}, errors.New("IAM用户信息缺少openid")
	}
	authIdentity, err := s.userRepo.QueryAuth(model.AuthTypeIAM, userInfo.OpenID)
	switch {
	case err == nil:
		user, err := s.userRepo.QueryById(authIdentity.UserID)
		return userToDomainUser(user), err
	case errors.Is(err, gorm.ErrRecordNotFound):
		return s.createIAMUser(userInfo)
	default:
		return model.DomainUser{}, err
	}
}

func (s *IAMOAuthService) createIAMUser(userInfo *oauth2Client.UserInfo) (model.DomainUser, error) {
	username := strings.TrimSpace(userInfo.Username)
	if username == "" {
		username = strings.TrimSpace(userInfo.OpenID)
	}
	username, err := s.availableUsername(username)
	if err != nil {
		return model.DomainUser{}, err
	}
	nickname := strings.TrimSpace(userInfo.DisplayName)
	if nickname == "" {
		nickname = username
	}

	user := model.User{
		Username:  username,
		Password:  "",
		Nickname:  nickname,
		LoginTime: time.Now(),
	}
	if err = s.userRepo.Insert(&user); err != nil {
		return model.DomainUser{}, err
	}
	if err = s.userRepo.CreateAuth(&model.AuthIdentity{
		UserID:   user.Id,
		Type:     model.AuthTypeIAM,
		Identity: userInfo.OpenID,
		Name:     username,
	}); err != nil {
		return model.DomainUser{}, err
	}
	return userToDomainUser(user), nil
}

func (s *IAMOAuthService) availableUsername(base string) (string, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "iam-user"
	}
	_, err := s.userRepo.QueryByUsername(base)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return base, nil
	case err == nil:
		return "", ErrUsernameDuplicate
	default:
		return "", err
	}
}

func (s *IAMOAuthService) saveOpenID(ctx context.Context, token *oauth2Client.TokenResponse) error {
	if token == nil || strings.TrimSpace(token.AccessToken) == "" || strings.TrimSpace(token.OpenID) == "" {
		return ErrOAuthTokenMissingOpenID
	}
	if s.redis == nil {
		return nil
	}
	ttl := time.Duration(token.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	if err := s.redis.Set(ctx, oauthOpenIDKeyPrefix+token.AccessToken, token.OpenID, ttl).Err(); err != nil {
		slog.Error("save oauth openid failed", "err", err)
		return err
	}
	return nil
}

func (s *IAMOAuthService) ensureEnabled() error {
	if !s.enabled {
		return ErrOAuthDisabled
	}
	if s.client == nil {
		return ErrOAuthClientNotReady
	}
	return nil
}

func getIAMOAuthScopes(oauthCfg config.OAuth) []string {
	if len(oauthCfg.Scopes) > 0 {
		return oauthCfg.Scopes
	}
	return strings.Fields(oauthCfg.Scope)
}

func userToDomainUser(user model.User) model.DomainUser {
	return model.DomainUser{
		Id:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		About:       user.About,
		CreatedTime: user.CreatedTime,
		UpdatedTime: user.UpdatedTime,
		LoginTime:   user.LoginTime,
	}
}

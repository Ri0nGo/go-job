package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-job/internal/model"
	"resty.dev/v3"
	"strconv"
	"time"
)

var (
	authorizationURL = "https://github.com/login/oauth/authorize"
	accessTokenURL   = "https://github.com/login/oauth/access_token"
	userInfoURL      = "https://api.github.com/user"
)

type OAuth2Service struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func NewOAuth2Service(clientID, clientSecret, redirectURL string) *OAuth2Service {
	return &OAuth2Service{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (o *OAuth2Service) GetAuthUrl(ctx context.Context, state string) string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s",
		authorizationURL, o.clientID, o.redirectURL, state)
}

func (o *OAuth2Service) GetAuthIdentity(ctx context.Context, code string) (model.AuthIdentity, error) {
	accessToken, err := o.getAccessToken(ctx, code)
	if err != nil {
		return model.AuthIdentity{}, err
	}
	return o.getUserInfo(ctx, accessToken)
}

func (o *OAuth2Service) getAccessToken(ctx context.Context, code string) (string, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	params := map[string]string{
		"client_id":     o.clientID,
		"client_secret": o.clientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	}
	resp, err := resty.New().SetContext(ctx).
		SetHeaders(header).
		SetTimeout(time.Second * 3).
		R().
		SetPathParams(params).
		Post(accessTokenURL)
	if err != nil {
		return "", err
	}

	var atResp accessTokenResp
	err = json.Unmarshal(resp.Bytes(), &atResp)
	if err != nil {
		return "", err
	}
	if atResp.Error != "" {
		return "", errors.New(atResp.Error)
	}
	return atResp.AccessToken, nil
}

func (o *OAuth2Service) getUserInfo(ctx context.Context, accessToken string) (model.AuthIdentity, error) {
	header := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Content-Type":  "application/json",
	}
	resp, err := resty.New().SetContext(ctx).
		SetHeaders(header).
		SetTimeout(time.Second * 3).
		R().
		Get(userInfoURL)
	if err != nil {
		return model.AuthIdentity{}, err
	}
	var userInfo userInfoResp
	if err = json.Unmarshal(resp.Bytes(), &userInfo); err != nil {
		return model.AuthIdentity{}, err
	}
	return model.AuthIdentity{
		Type:     model.AuthTypeGithub,
		Identity: strconv.Itoa(userInfo.Id),
		Name:     userInfo.Login,
	}, nil
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

type userInfoResp struct {
	Login     string `json:"login"`      // username
	Id        int    `json:"id"`         // user id
	NodeId    string `json:"node_id"`    // 节点ID
	AvatarUrl string `json:"avatar_url"` // 头像地址
	HtmlUrl   string `json:"html_url"`   // 前端个人主页地址
}

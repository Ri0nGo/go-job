package qq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-job/internal/model"
	"net/url"
	"resty.dev/v3"
	"time"
)

// 参考：https://wiki.connect.qq.com/%e5%bc%80%e5%8f%91%e6%94%bb%e7%95%a5_server-side
// https://wiki.connect.qq.com/%e4%bd%bf%e7%94%a8authorization_code%e8%8e%b7%e5%8f%96access_token

var (
	authorizationURL = "https://graph.qq.com/oauth2.0/authorize"
	accessTokenURL   = "https://graph.qq.com/oauth2.0/token"
	openIdUrl        = "https://graph.qq.com/oauth2.0/me"
	userInfoURL      = "https://graph.qq.com/user/get_user_info"
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
	return fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		authorizationURL, o.clientID, o.redirectURL, state)
}

func (o *OAuth2Service) GetAuthIdentity(ctx context.Context, code string) (model.AuthIdentity, error) {
	// todo 待实现
	accessToken, err := o.getAccessToken(ctx, code)
	if err != nil {
		return model.AuthIdentity{}, err
	}
	openId, err := o.getOpenId(ctx, accessToken)
	if err != nil {
		return model.AuthIdentity{}, err
	}

	return o.getUserInfo(ctx, accessToken, openId)
}

func (o *OAuth2Service) getAccessToken(ctx context.Context, code string) (string, error) {
	params := map[string]string{
		"client_id":     o.clientID,
		"client_secret": o.clientSecret,
		"code":          code,
		"redirect_uri":  url.PathEscape(o.redirectURL),
		"grant_type":    "authorization_code",
		"fmt":           "json",
	}

	// 发起请求
	resp, err := resty.New().SetContext(ctx).
		R().
		SetQueryParams(params).
		SetTimeout(3 * time.Second).
		Get(accessTokenURL)

	if err != nil {
		return "", err
	}

	// 解析响应 JSON
	var result accessTokenResp
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", errors.New(result.Error)
	}

	return result.AccessToken, nil
}

func (o *OAuth2Service) getOpenId(ctx context.Context, accessToken string) (string, error) {
	params := map[string]string{
		"access_token": accessToken,
		"fmt":          "json",
	}
	resp, err := resty.New().
		R().
		SetQueryParams(params).
		Get(openIdUrl)
	if err != nil {
		return "", err
	}
	var openResp openIdResp
	if err = json.Unmarshal(resp.Bytes(), &openResp); err != nil {
		return "", err
	}
	return openResp.OpenId, nil
}

func (o *OAuth2Service) getUserInfo(ctx context.Context, accessToken, openId string) (model.AuthIdentity, error) {
	header := map[string]string{
		"access_token":       accessToken,
		"oauth_consumer_key": o.clientID,
		"openid":             openId,
		"fmt":                "json",
	}
	resp, err := resty.New().
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
	if userInfo.Ret != 0 {
		return model.AuthIdentity{}, errors.New(userInfo.Msg)
	}
	return model.AuthIdentity{
		Type:     model.AuthTypeQQ,
		Identity: openId,
		Name:     userInfo.Nickname,
	}, nil
}

type accessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Error        string `json:"error"`
}

type openIdResp struct {
	ClientId string `json:"client_id"`
	OpenId   string `json:"openid"`
}

type userInfoResp struct {
	Ret      int    `json:"ret"` // 状态码
	Msg      string `json:"msg"`
	IsLost   string `json:"is_lost"`
	Nickname string `json:"nickname"`
}

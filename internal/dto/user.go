package dto

import "go-job/internal/model"

type ReqUserEmailBind struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ReqOAuth2Bind struct {
	Code     string `json:"code"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ReqOAuth2UnBind struct {
	AuthType model.AuthType `json:"auth_type"`
}

type RespUserSecurity struct {
	Email  string `json:"email"`
	QQ     bool   `json:"qq"`
	Github bool   `json:"github"`
}

type RespOAuth2Code struct {
	Err          string `json:"err"`
	RedirectPage string `json:"redirect_page"`
	Platform     string `json:"platform"`
	Token        string `json:"-"` // 默认放在Header中
	UID          int    `json:"-"`
}

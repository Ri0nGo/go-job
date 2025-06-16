package dto

type ReqUserEmailBind struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ReqOAuth2Bind struct {
	Key      string `json:"key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespUserSecurity struct {
	Email  string `json:"email"`
	QQ     bool   `json:"qq"`
	Github bool   `json:"github"`
}

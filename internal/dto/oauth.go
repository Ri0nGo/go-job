package dto

type OAuthInfo struct {
	Enabled      bool   `json:"enabled"`
	AuthURL      string `json:"auth_url"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
}

type ReqOAuthLogin struct {
	Code string `json:"code" binding:"required"`
}

type RespOAuthLogin struct {
	ID int `json:"id"`
}

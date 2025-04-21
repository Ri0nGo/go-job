package dto

type ReqUserEmailBind struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

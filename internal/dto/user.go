package dto

import "go-job/internal/model"

type ReqUserBind struct {
	Type  model.UserBindType `json:"type" binding:"required"`
	Email string             `json:"email"`
	Code  string             `json:"code"`
}

type ReqCodeSend struct {
	Type  model.UserBindType `json:"type" binding:"required"`
	Email string             `json:"email"`
}

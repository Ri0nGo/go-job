package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Id              int       `json:"id" gorm:"primary_key"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Nickname        string    `json:"nickname"`
	About           string    `json:"about"`
	CreatedTime     time.Time `json:"created_at" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime     time.Time `json:"updated_at" gorm:"column:updated_time;autoUpdateTime"`
	ConfirmPassword string    `json:"confirm_password" gorm:"-"`
}

func (u *User) TableName() string {
	return "user"
}

type DomainUser struct {
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	About       string    `json:"about"`
	CreatedTime time.Time `json:"created_time"`
}

func (u *DomainUser) TableName() string {
	return "user"
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int `json:"uid"`
}

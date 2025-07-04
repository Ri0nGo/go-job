package model

import (
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type AuthType uint8

func (a AuthType) String() string {
	switch a {
	case AuthTypeGithub:
		return "github"
	case AuthTypeQQ:
		return "qq"
	default:
		return "unknown(" + strconv.Itoa(int(a)) + ")"
	}
}

const (
	AuthTypeGithub AuthType = iota + 1
	AuthTypeQQ
)

type Auth2Scene string

const (
	Auth2SceneLoginPage           Auth2Scene = "login_page"
	Auth2SceneAccountSecurityPage            = "account_page"
)

type User struct {
	Id              int       `json:"id" gorm:"primary_key"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Nickname        string    `json:"nickname"`
	About           string    `json:"about"`
	Email           *string   `json:"email"`
	CreatedTime     time.Time `json:"created_time" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime     time.Time `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"`
	LoginTime       time.Time `json:"login_time" gorm:"column:login_time"`
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
	Email       string    `json:"email"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
	LoginTime   time.Time `json:"login_time"`
}

func (u *DomainUser) TableName() string {
	return "user"
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int `json:"uid"`
}

type AuthIdentity struct {
	ID          int       `json:"id" gorm:"primary_key"` // 自增id
	UserID      int       `json:"user_id"`               // 用户ID， 一对多关系
	Type        AuthType  `json:"type"`                  // 授权类型 1: github
	Identity    string    `json:"identity"`              // 授权唯一标识
	Name        string    `json:"name"`                  // 授权平台的用户名
	CreatedTime time.Time `json:"created_at" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time `json:"updated_at" gorm:"column:updated_time;autoUpdateTime"`
}

func (a *AuthIdentity) TableName() string {
	return "auth_identity"
}

type UserAuthInfo struct {
	ID       int            `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Auths    []AuthIdentity `json:"auths"`
}

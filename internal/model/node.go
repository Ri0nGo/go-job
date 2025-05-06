package model

import "time"

type Node struct {
	Id          int       `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Address     string    `json:"address" binding:"required"`
	CreatedTime time.Time `json:"created_at" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time `json:"updated_at" gorm:"column:updated_time;autoUpdateTime"`
	Online      bool      `json:"online" gorm:"-"`
	CheckTime   time.Time `json:"check_time" gorm:"-"`
}

func (Node) TableName() string {
	return "node"
}

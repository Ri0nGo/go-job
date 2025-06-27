package model

import "time"

type SystemOperationLog struct {
	Id          int       `json:"id" gorm:"primary_key"`
	UserId      int       `json:"user_id" gorm:"column:user_id"`
	ClientIP    string    `json:"client_ip" gorm:"column:client_ip"`
	Method      string    `json:"method" gorm:"column:method"`
	UA          string    `json:"ua" gorm:"column:ua"`
	URL         string    `json:"url"`
	StatusCode  int       `json:"status_code" gorm:"column:status_code"`
	Response    string    `json:"response" gorm:"column:response"`
	Request     string    `json:"request"`
	Title       string    `json:"title" gorm:"column:title"`
	CreatedTime time.Time `json:"created_time" gorm:"column:created_time;autoCreateTime"`
}

func (s *SystemOperationLog) TableName() string {
	return "sys_opt_log"
}

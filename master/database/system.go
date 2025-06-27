package database

import "go-job/internal/model"

func CreateSysOptLogWithMySQL(log model.SystemOperationLog) error {
	return mysqlDb.Create(&log).Error
}

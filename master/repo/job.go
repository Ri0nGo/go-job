package repo

import (
	"gorm.io/gorm"
)

type IJobRepo interface {
}

type JobRepo struct {
	mysqlDB *gorm.DB
}

func NewJobRepo(mysqlDB *gorm.DB) IJobRepo {
	return &JobRepo{
		mysqlDB: mysqlDB,
	}
}

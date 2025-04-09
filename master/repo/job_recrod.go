package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type IJobRecordRepo interface {
	QueryById(id int) (model.JobRecord, error)
	Inserts([]model.JobRecord) error
	Insert(*model.JobRecord) error
	Delete(id int) error
	QueryList(page model.Page) (model.Page, error)
}

type JobRecordRepo struct {
	mysqlDB *gorm.DB
}

func (j *JobRecordRepo) QueryById(id int) (model.JobRecord, error) {
	var job model.JobRecord
	err := j.mysqlDB.First(&job, id).Error
	return job, err
}

func (j *JobRecordRepo) Inserts(jobs []model.JobRecord) error {
	if len(jobs) == 0 {
		return nil
	}
	return j.mysqlDB.Create(&jobs).Error
}

func (j *JobRecordRepo) Insert(job *model.JobRecord) error {
	return j.mysqlDB.Create(job).Error
}

func (j *JobRecordRepo) Delete(id int) error {
	return j.mysqlDB.Where("id = ?", id).Delete(&model.JobRecord{}).Error
}

func (j *JobRecordRepo) QueryList(page model.Page) (model.Page, error) {
	return paginate.PaginateList[model.JobRecord](j.mysqlDB, page.PageNum, page.PageSize)
}

func NewJobRecordRepo(mysqlDB *gorm.DB) IJobRecordRepo {
	return &JobRecordRepo{
		mysqlDB: mysqlDB,
	}
}

package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type IJobRepo interface {
	QueryById(id int) (model.Job, error)
	Inserts([]model.Job) error
	Update(model.Job) error
	Delete(id int) error
	QueryList(page model.Page) (model.Page, error)
}

type JobRepo struct {
	mysqlDB *gorm.DB
}

func (j *JobRepo) QueryById(id int) (model.Job, error) {
	var job model.Job
	err := j.mysqlDB.First(&job, id).Error
	return job, err
}

func (j *JobRepo) Inserts(jobs []model.Job) error {
	if len(jobs) == 0 {
		return nil
	}
	return j.mysqlDB.Create(&jobs).Error
}

func (j *JobRepo) Update(job model.Job) error {
	return j.mysqlDB.Where("id = ?", job.Id).Updates(&job).Error
}

func (j *JobRepo) Delete(id int) error {
	return j.mysqlDB.Where("id = ?", id).Delete(&model.Job{}).Error
}

func (j *JobRepo) QueryList(page model.Page) (model.Page, error) {
	return paginate.PaginateList[model.Job](j.mysqlDB, page.PageNum, page.PageSize)
}

func NewJobRepo(mysqlDB *gorm.DB) IJobRepo {
	return &JobRepo{
		mysqlDB: mysqlDB,
	}
}

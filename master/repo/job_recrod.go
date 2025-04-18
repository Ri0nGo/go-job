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
	QueryList(page model.Page, jobId int) (model.Page, error)
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

func (j *JobRecordRepo) QueryList(page model.Page, jobId int) (model.Page, error) {
	// TODO 这里后面还需要优化查询方式，感觉还需要对分页查询做封装
	return paginate.PaginateListV2[model.JobRecordSummary](j.mysqlDB, page.PageNum, page.PageSize, func(db *gorm.DB) *gorm.DB {
		return db.Where("job_id = ?", jobId).Order("id desc")
	})
}

func NewJobRecordRepo(mysqlDB *gorm.DB) IJobRecordRepo {
	return &JobRecordRepo{
		mysqlDB: mysqlDB,
	}
}

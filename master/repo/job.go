package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type IJobRepo interface {
	QueryById(id int) (model.Job, error)
	QueryByIds(id []int) ([]model.Job, error)
	QueryByNodeId(nodeId int) ([]model.Job, error)
	Insert(*model.Job) error
	Inserts([]model.Job) error
	Update(*model.Job) error
	Delete(id int) error
	QueryListByUID(uid int, page model.Page) (model.Page, error)
	QuerySummary(uid int) ([]model.JobStatusCount, error)
}

type JobRepo struct {
	mysqlDB *gorm.DB
}

func (j *JobRepo) QueryById(id int) (model.Job, error) {
	var job model.Job
	err := j.mysqlDB.First(&job, id).Error
	return job, err
}

func (j *JobRepo) QueryByIds(ids []int) ([]model.Job, error) {
	var jobs []model.Job
	err := j.mysqlDB.Where("id IN ?", ids).Find(&jobs).Error
	return jobs, err
}

func (j *JobRepo) Inserts(jobs []model.Job) error {
	if len(jobs) == 0 {
		return nil
	}
	return j.mysqlDB.Create(&jobs).Error
}

func (j *JobRepo) Insert(job *model.Job) error {
	return j.mysqlDB.Create(job).Error
}

func (j *JobRepo) Update(job *model.Job) error {
	if job.Id == 0 {
		return ErrorIDIsZero
	}
	return j.mysqlDB.Updates(job).Error
}

func (j *JobRepo) Delete(id int) error {
	return j.mysqlDB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ?", id).Delete(&model.Job{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		return tx.Where("job_id = ?", id).Delete(&model.JobRecord{}).Error
	})
}

func (j *JobRepo) QueryByNodeId(nodeId int) ([]model.Job, error) {
	var jobs []model.Job
	err := j.mysqlDB.Where("node_id = ?", nodeId).Find(&jobs).Error
	return jobs, err
}

func (j *JobRepo) QueryListByUID(uid int, page model.Page) (model.Page, error) {
	return paginate.PaginateListV2[model.Job](j.mysqlDB, page, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", uid)
	})
}

func (j *JobRepo) QuerySummary(uid int) ([]model.JobStatusCount, error) {
	var data []model.JobStatusCount
	err := j.mysqlDB.Model(&model.Job{}).Select("active, COUNT(*) as count").
		Where("user_id = ?", uid).
		Group("active").
		Scan(&data).Error
	return data, err
}

func NewJobRepo(mysqlDB *gorm.DB) IJobRepo {
	return &JobRepo{
		mysqlDB: mysqlDB,
	}
}

package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
	"time"
)

type IJobRecordRepo interface {
	QueryById(id int) (model.JobRecord, error)
	Inserts([]model.JobRecord) error
	Insert(*model.JobRecord) error
	Delete(id int) error
	QueryList(page model.Page, jobId int) (model.Page, error)
	QueryLastList(page model.Page) (model.Page, error)
	QueryDayStatus(begin, end time.Time) ([]model.JobRecordDayStatusCount, error)
	QueryJobStatus(begin, end time.Time) ([]model.JobRecordJobStatusCount, error)
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
	return paginate.PaginateListV2[model.JobRecordSummary](j.mysqlDB, page, func(db *gorm.DB) *gorm.DB {
		return db.Where("job_id = ?", jobId)
	})
}

func (j *JobRecordRepo) QueryDayStatus(being, end time.Time) ([]model.JobRecordDayStatusCount, error) {
	var jobs []model.JobRecordDayStatusCount
	err := j.mysqlDB.Model(&model.JobRecord{}).Select("DATE(start_time) AS date, status, COUNT(*) as count").
		Where("start_time >= ? AND end_time <= ?", being, end).
		Group("date, status").
		Order("date, status").
		Find(&jobs).Error
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (j *JobRecordRepo) QueryJobStatus(being, end time.Time) ([]model.JobRecordJobStatusCount, error) {
	var jobs []model.JobRecordJobStatusCount
	err := j.mysqlDB.Model(&model.JobRecord{}).Select("job_id, status, COUNT(*) as count").
		Where("start_time >= ? AND end_time <= ?", being, end).
		Group("job_id, status").
		Order("job_id, status").
		Find(&jobs).Error
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (j *JobRecordRepo) QueryLastList(page model.Page) (model.Page, error) {
	var jobs []model.JobLastRecord
	err := j.mysqlDB.Table("job_record AS r").
		Select("r.id, j.id AS job_id, j.name AS job_name, n.id AS node_id, n.name AS node_name, r.start_time, r.end_time, r.status").
		Joins("JOIN job j ON r.job_id = j.id").
		Joins("JOIN node n ON n.id = j.node_id").
		Order("r.start_time DESC, r.id DESC").
		Limit(page.PageSize).
		Scan(&jobs).Error
	if err != nil {
		return page, err
	}
	page.Total = int64(len(jobs))
	page.Data = jobs
	return page, nil
}

func NewJobRecordRepo(mysqlDB *gorm.DB) IJobRecordRepo {
	return &JobRecordRepo{
		mysqlDB: mysqlDB,
	}
}

package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type INodeRepo interface {
	QueryById(id int) (model.Node, error)
	Inserts([]model.Node) error
	Update(model.Node) error
	Delete(id int) error
	QueryList(page model.Page) (model.Page, error)
}

type NodeRepo struct {
	mysqlDB *gorm.DB
}

func (j *NodeRepo) QueryById(id int) (model.Node, error) {
	var job model.Node
	err := j.mysqlDB.First(&job, id).Error
	return job, err
}

func (j *NodeRepo) Inserts(jobs []model.Node) error {
	if len(jobs) == 0 {
		return nil
	}
	return j.mysqlDB.Create(&jobs).Error
}

func (j *NodeRepo) Update(job model.Node) error {
	return j.mysqlDB.Updates(&job).Error
}

func (j *NodeRepo) Delete(id int) error {
	return j.mysqlDB.Where("id = ?", id).Delete(&model.Node{}).Error
}

func (j *NodeRepo) QueryList(page model.Page) (model.Page, error) {
	return paginate.PaginateList[model.Node](j.mysqlDB, page.PageNum, page.PageSize)
}

func NewNodeRepo(mysqlDB *gorm.DB) INodeRepo {
	return &NodeRepo{
		mysqlDB: mysqlDB,
	}
}

package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type INodeRepo interface {
	QueryById(id int) (model.Node, error)
	QueryByIds(ids []int) ([]model.Node, error)
	Inserts([]model.Node) error
	Update(model.Node) error
	Delete(id int) error
	QueryList(page model.Page) (model.Page, error)
}

type NodeRepo struct {
	mysqlDB *gorm.DB
}

func (j *NodeRepo) QueryById(id int) (model.Node, error) {
	var node model.Node
	err := j.mysqlDB.First(&node, id).Error
	return node, err
}

func (j *NodeRepo) QueryByIds(ids []int) ([]model.Node, error) {
	var nodes []model.Node
	err := j.mysqlDB.Where("id IN (?)", ids).Find(&nodes).Error
	return nodes, err
}

func (j *NodeRepo) Inserts(nodes []model.Node) error {
	if len(nodes) == 0 {
		return nil
	}
	return j.mysqlDB.Create(&nodes).Error
}

func (j *NodeRepo) Update(node model.Node) error {
	if node.Id == 0 {
		return ErrorIDIsZero
	}
	return j.mysqlDB.Updates(&node).Error
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

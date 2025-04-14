package service

import (
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/repo"
)

type INodeService interface {
	GetNode(id int) (model.Node, error)
	GetNodeList(page model.Page) (model.Page, error)
	AddNode(job model.Node) error
	DeleteNode(id int) error
	UpdateNode(job model.Node) error
}

type NodeService struct {
	NodeRepo repo.INodeRepo
	JobRepo  repo.IJobRepo
}

func (s *NodeService) GetNode(id int) (model.Node, error) {
	return s.NodeRepo.QueryById(id)
}

func (s *NodeService) GetNodeList(page model.Page) (model.Page, error) {
	return s.NodeRepo.QueryList(page)
}

func (s *NodeService) AddNode(node model.Node) error {
	if !utils.IsValidIPv4Address(node.Address) {
		return ErrInvalidAddress
	}
	return s.NodeRepo.Inserts([]model.Node{node})
}

func (s *NodeService) DeleteNode(id int) error {
	nodes, err := s.JobRepo.QueryByNodeId(id)
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return ErrJobUseCurrentNode
	}
	return s.NodeRepo.Delete(id)
}

func (s *NodeService) UpdateNode(node model.Node) error {
	if !utils.IsValidIPv4Address(node.Address) {
		return ErrInvalidAddress
	}
	return s.NodeRepo.Update(node)
}

func NewNodeService(nodeRepo repo.INodeRepo, jobRepo repo.IJobRepo) INodeService {
	return &NodeService{
		NodeRepo: nodeRepo,
		JobRepo:  jobRepo,
	}
}

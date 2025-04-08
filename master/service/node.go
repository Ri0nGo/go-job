package service

import (
	"errors"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/repo"
)

var errInvalidAddress = errors.New("invalid address")

type INodeService interface {
	GetNode(id int) (model.Node, error)
	GetNodeList(page model.Page) (model.Page, error)
	AddNode(job model.Node) error
	DeleteNode(id int) error
	UpdateNode(job model.Node) error
}

type NodeService struct {
	NodeRepo repo.INodeRepo
}

func (s *NodeService) GetNode(id int) (model.Node, error) {
	return s.NodeRepo.QueryById(id)
}

func (s *NodeService) GetNodeList(page model.Page) (model.Page, error) {
	return s.NodeRepo.QueryList(page)
}

func (s *NodeService) AddNode(node model.Node) error {
	if !utils.IsValidIPv4Address(node.Address) {
		return errInvalidAddress
	}
	return s.NodeRepo.Inserts([]model.Node{node})
}

func (s *NodeService) DeleteNode(id int) error {
	return s.NodeRepo.Delete(id)
}

func (s *NodeService) UpdateNode(node model.Node) error {
	if !utils.IsValidIPv4Address(node.Address) {
		return errInvalidAddress
	}
	return s.NodeRepo.Update(node)
}

func NewNodeService(nodeRepo repo.INodeRepo) INodeService {
	return &NodeService{
		NodeRepo: nodeRepo,
	}
}

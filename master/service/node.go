package service

import (
	"go-job/internal/model"
	"go-job/master/pkg/metrics"
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
	data, err := s.NodeRepo.QueryList(page)
	if err != nil {
		return data, err
	}

	// 处理节点指标数据
	nodes := data.Data.([]model.Node)
	nodeMetric := metrics.GetNodeMetrics()
	for i, node := range nodes {
		if m, ok := nodeMetric.Get(node.Id); ok {
			nodes[i].Online = m.Online
			nodes[i].CheckTime = m.CheckTime
		}
	}
	return data, nil
}

func (s *NodeService) AddNode(node model.Node) error {
	if err := s.NodeRepo.Insert(&node); err != nil {
		return err
	}
	metrics.GetNodeMetrics().SetAndCheck(node.Id, node)
	return nil
}

func (s *NodeService) DeleteNode(id int) error {
	nodes, err := s.JobRepo.QueryByNodeId(id)
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return ErrJobUseCurrentNode
	}
	if err := s.NodeRepo.Delete(id); err != nil {
		return err
	}
	metrics.GetNodeMetrics().Remove(id)
	return nil
}

func (s *NodeService) UpdateNode(node model.Node) error {
	if err := s.NodeRepo.Update(node); err != nil {
		return err
	}
	metrics.GetNodeMetrics().Remove(node.Id)
	metrics.GetNodeMetrics().Set(node.Id, node)
	return nil
}

func NewNodeService(nodeRepo repo.INodeRepo, jobRepo repo.IJobRepo) INodeService {
	return &NodeService{
		NodeRepo: nodeRepo,
		JobRepo:  jobRepo,
	}
}

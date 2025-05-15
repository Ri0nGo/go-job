package service

import (
	"context"
	"errors"
	"fmt"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/httpClient"
	"go-job/master/pkg/metrics"
	"go-job/master/repo"
	"log/slog"
	"time"
)

type INodeService interface {
	GetNode(id int) (model.Node, error)
	GetNodeList(page model.Page) (model.Page, error)
	AddNode(job model.Node) error
	DeleteNode(id int) error
	UpdateNode(job model.Node) error
	InstallRef(req dto.ReqNodeRef) (any, error)
	NodeInfo(id int) (any, error)
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

func (s *NodeService) InstallRef(req dto.ReqNodeRef) (any, error) {
	node, err := s.NodeRepo.QueryById(req.Id)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/api/go-job/node/install_ref", node.Address)
	resp, err := httpClient.PostJson(context.Background(), url, nil, req, time.Second*10)
	if err != nil {
		slog.Error("send pkg name to node error by install ref", "url", url,
			"req", req, "err", err)
		return nil, err
	}

	nodeResp, err := httpClient.ParseResponse(resp)
	if err != nil {
		slog.Error("send install pkg to node error", "url", url, "resp", resp, "err", err)
		return nil, err
	}
	if nodeResp.Code != 0 {
		slog.Error("resp code isn't zero", "resp", resp)
		return nil, errors.New("resp code isn't zero in install python pkg name")
	}
	return nodeResp.Data, nil
}

func (s *NodeService) NodeInfo(id int) (any, error) {
	node, err := s.NodeRepo.QueryById(id)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s/api/go-job/node/info", node.Address)
	resp, err := httpClient.GetJson(context.Background(), url, nil, nil, time.Second*10)
	if err != nil {
		slog.Error("send pkg name to node error by install ref", "url", url, "err", err)
		return nil, err
	}

	nodeResp, err := httpClient.ParseResponse(resp)
	if err != nil {
		slog.Error("send install pkg to node error", "url", url, "resp", resp, "err", err)
		return nil, err
	}
	if nodeResp.Code != 0 {
		slog.Error("resp code isn't zero", "resp", resp)
		return nil, errors.New("resp code isn't zero in install python pkg name")
	}
	return nodeResp.Data, nil
}

func NewNodeService(nodeRepo repo.INodeRepo, jobRepo repo.IJobRepo) INodeService {
	return &NodeService{
		NodeRepo: nodeRepo,
		JobRepo:  jobRepo,
	}
}

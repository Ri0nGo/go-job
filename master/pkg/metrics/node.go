package metrics

import (
	"context"
	"go-job/internal/model"
	"log/slog"
	"net"
	"sync"
	"time"
)

// note 存在的问题
// 1. 没有提供注册自定义检查方法的机制
// 2. 所有节点使用相同的检查逻辑和超时设置，不支持节点级别的配置
// 3. 缺少错误处理机制，特别是在网络检测过程中
// 4. 如果节点数量很大，当前的实现可能会导致性能问题

var (
	nodeMetricsInstance        *NodeMetrics
	onceNode                   sync.Once
	defaultCheckoutNodeTimeout = 2 * time.Second
	defaultInterval            = 30 * time.Second
)

type NodeOption func(m *NodeMetrics)

func WithTimeout(timeout time.Duration) NodeOption {
	return func(m *NodeMetrics) {
		m.timeout = timeout
	}
}

func WithInterval(interval time.Duration) NodeOption {
	return func(m *NodeMetrics) {
		m.interval = interval
	}
}

type NodeMetrics struct {
	mux      sync.RWMutex
	ctx      context.Context
	nodes    map[int]*NodeMetric
	timeout  time.Duration // 检测存活超时时间
	interval time.Duration // 检测节点间隔
}

type NodeMetric struct {
	model.Node
}

func (m *NodeMetrics) Set(nodeId int, node model.Node) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.nodes[nodeId] = &NodeMetric{
		Node: node,
	}
}

func (m *NodeMetrics) Get(nodeId int) (*NodeMetric, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	result, ok := m.nodes[nodeId]
	return result, ok
}

func (m *NodeMetrics) Remove(nodeId int) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.nodes, nodeId)
}

func (m *NodeMetrics) BuildNodeMetric(nodes map[int]model.Node) *NodeMetrics {
	m.mux.Lock()
	defer m.mux.Unlock()

	for nodeId, node := range nodes {
		m.nodes[nodeId] = &NodeMetric{
			Node: node,
		}
	}
	return m
}

func isConnected(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		slog.Error("check node error", "addr", addr, "err", err)
		return false
	}
	conn.Close()
	return true
}

func (m *NodeMetrics) checkNodesMetric() {
	m.mux.RLock()
	addrs := make(map[int]string)
	for id, metric := range m.nodes {
		addrs[id] = metric.Address
	}
	m.mux.RUnlock()

	results := make(map[int]bool)
	for id, addr := range addrs {
		results[id] = isConnected(addr, m.timeout)
	}

	m.mux.Lock()
	for id, status := range results {
		if nm, ok := m.nodes[id]; ok {
			nm.Online = status
			nm.CheckTime = time.Now()
		}
	}
	m.mux.Unlock()
}

func (m *NodeMetrics) Monitor() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.checkNodesMetric()
	for {
		select {
		case <-ticker.C:
			m.checkNodesMetric()
		case <-m.ctx.Done():
			return

		}
	}
}

func newNodeMetrics(ctx context.Context) *NodeMetrics {
	return &NodeMetrics{
		ctx:      ctx,
		nodes:    make(map[int]*NodeMetric),
		timeout:  defaultCheckoutNodeTimeout,
		interval: defaultInterval,
	}
}

func InitNodeMetrics(ctx context.Context, nodes map[int]model.Node, opts ...NodeOption) {
	onceNode.Do(func() {
		nodeMetricsInstance = newNodeMetrics(ctx).BuildNodeMetric(nodes)
		for _, opt := range opts {
			opt(nodeMetricsInstance)
		}
	})
}

func GetNodeMetrics() *NodeMetrics {
	return nodeMetricsInstance
}

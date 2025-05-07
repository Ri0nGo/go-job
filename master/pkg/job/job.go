package job

import (
	"context"
	"go-job/internal/model"
	"go-job/master/pkg/config"
	"go-job/master/pkg/metrics"
	"go-job/master/pkg/notify"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
)

func InitGlobalData(mysqlDB *gorm.DB, jobSvc service.IJobService,
	notifyStore notify.INotifyStore) error {
	// 查询所有的job
	jobs, err := queryAllJobs(mysqlDB)
	if err != nil {
		panic(err)
	}

	nodeM, err := queryAllNodes(mysqlDB)
	if err != nil {
		panic(err)
	}

	initJobData(nodeM, jobs, jobSvc, notifyStore)

	initNodeMetrics(nodeM)

	return nil
}

func initJobData(nodeM map[int]model.Node, jobs []model.Job,
	jobSvc service.IJobService, notifyStore notify.INotifyStore) {

	for _, job := range jobs {
		// 发送job到node
		node := nodeM[job.NodeID]
		err := jobSvc.SendJobToNode(job, node, service.SendJobByCreate)
		if err != nil {
			slog.Error("init job to node error", "job name",
				job.Name, "job id", job.Id, "err", err)
		}

		// 初始化任务通知数据
		if job.Internal.Notify.NotifyStatus == model.NotifyStatusEnabled {
			notifyStore.Set(context.Background(), job.Id, service.GenNotifyConfig(job))
		}
	}
}

func initNodeMetrics(nodeM map[int]model.Node) {
	metrics.InitNodeMetrics(context.Background(), nodeM,
		metrics.WithNodeTimeout(config.App.Metrics.Node.Timeout),
		metrics.WithNodeInterval(config.App.Metrics.Node.Interval))
	go metrics.GetNodeMetrics().Monitor()
}

func queryAllJobs(db *gorm.DB) ([]model.Job, error) {
	var jobs []model.Job
	err := db.Find(&jobs).Error
	return jobs, err
}

func queryAllNodes(db *gorm.DB) (map[int]model.Node, error) {
	var (
		nodes   []model.Node
		nodeMap = make(map[int]model.Node)
	)
	err := db.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		nodeMap[node.Id] = node

	}
	return nodeMap, err
}

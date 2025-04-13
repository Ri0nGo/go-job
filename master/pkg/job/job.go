package job

import (
	"go-job/internal/model"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
)

func InitJobDataToNode(mysqlDB *gorm.DB, jobSvc service.IJobService) error {
	// 查询所有的job
	jobs, err := queryAllJobs(mysqlDB)
	if err != nil {
		panic(err)
	}

	nodeM, err := queryAllNodes(mysqlDB)
	if err != nil {
		panic(err)
	}

	// 发送job到node
	for _, job := range jobs {
		//fmt.Println("send job to node info", job.Id, job.Name, job.CronExpr)
		node := nodeM[job.NodeID]
		err := jobSvc.SendJobToNode(job, node, service.SendJobByCreate)
		if err != nil {
			slog.Error("init job to node error", "job name",
				job.Name, "job id", job.Id, "err", err)
		}
	}
	return nil
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

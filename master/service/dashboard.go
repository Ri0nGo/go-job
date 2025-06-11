package service

import (
	"errors"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/metrics"
	"go-job/master/repo"
	"time"
)

var (
	ErrUnSupportChartKey = errors.New("无效的查询key")
)

type IDashboardService interface {
	GetDataSummary(uid int) (dto.RespDashboardTotalData, error) // 卡片头部数据
	GetChartData(req model.ReqDashboardChart, uid int) (map[string][]int, error)
}

type DashboardService struct {
	JobRepo       repo.IJobRepo
	JobRecordRepo repo.IJobRecordRepo
}

func (s *DashboardService) GetDataSummary(uid int) (dto.RespDashboardTotalData, error) {
	var (
		data      dto.RespDashboardTotalData
		onlineCnt int
	)

	// 获取节点数据
	nodes := metrics.GetNodeMetrics().All()
	totalNode := len(nodes)
	for _, node := range nodes {
		if node.Online {
			onlineCnt++
		}
	}
	data.Node = dto.NodeDashboard{
		Total:   totalNode,
		Online:  onlineCnt,
		Offline: totalNode - onlineCnt,
	}

	// 查询任务数量
	jobStatus, err := s.JobRepo.QuerySummary(uid)
	if err != nil {
		return data, err
	}
	var (
		inActiveCnt int
		activeCnt   int
	)
	for _, js := range jobStatus {
		if js.Active == model.JobStart {
			activeCnt = js.Count
		} else {
			inActiveCnt = js.Count
		}
	}
	data.Job = dto.JobDashboard{
		Total:    activeCnt + inActiveCnt,
		Active:   activeCnt,
		InActive: inActiveCnt,
	}

	return data, nil
}

func (s *DashboardService) GetChartData(req model.ReqDashboardChart, uid int) (map[string][]int, error) {
	beginTime := utils.TimestampToTime(req.Begin)
	endTime := utils.TimestampToTime(req.End)
	var result = make(map[string][]int)

	switch req.Key {
	case model.DashboardKeyJobStatus:
		data, err := s.JobRecordRepo.QueryJobStatusByUid(
			utils.TimestampToTime(req.Begin),
			utils.TimestampToTime(req.End),
			uid)
		if err != nil {
			return result, err
		}
		if len(data) == 0 {
			return result, nil
		}

		jobIds := make([]int, 0)
		for _, item := range data {
			jobIds = append(jobIds, item.JobId)
		}
		jobs, err := s.JobRepo.QueryByIds(jobIds)
		if err != nil {
			return result, err
		}
		jobMap := make(map[int]string)
		for _, job := range jobs {
			jobMap[job.Id] = job.Name
		}

		for _, name := range jobMap {
			result[name] = []int{0, 0, 0, 0}
		}
		for _, item := range data {
			if name, ok := jobMap[item.JobId]; ok {
				result[name][item.Status] = item.Count
			}
		}

	case model.DashboardKeyDayStatus:
		for beginTime.Before(endTime) {
			result[beginTime.Format(time.DateOnly)] = []int{0, 0, 0, 0}
			beginTime = beginTime.AddDate(0, 0, 1)
		}
		data, err := s.JobRecordRepo.QueryDayStatusByUid(
			utils.TimestampToTime(req.Begin),
			utils.TimestampToTime(req.End),
			uid)
		if err != nil {
			return result, err
		}
		for _, item := range data {
			result[item.Date.Format(time.DateOnly)][item.Status] = item.Count
		}

	default:
		return result, ErrUnSupportChartKey

	}
	return result, nil
}

func NewDashboardService(jobRepo repo.IJobRepo, jobRecordRepo repo.IJobRecordRepo) IDashboardService {
	return &DashboardService{
		JobRepo:       jobRepo,
		JobRecordRepo: jobRecordRepo,
	}
}

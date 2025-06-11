package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/master/service"
	"log/slog"
)

type DashboardApi struct {
	dashboardService service.IDashboardService
}

func NewDashboardApi(DashboardService service.IDashboardService) *DashboardApi {
	return &DashboardApi{
		dashboardService: DashboardService,
	}
}

// RegisterRoutes 注册节点模块路由
func (a *DashboardApi) RegisterRoutes(group *gin.RouterGroup) {
	DashboardGroup := group.Group("/dashboards")
	{
		DashboardGroup.GET("/summary", a.GetDashboardSummary)
		DashboardGroup.POST("/chart", a.GetDashboardChart)
	}
}

// GetDashboard 查询节点
func (a *DashboardApi) GetDashboardSummary(ctx *gin.Context) {
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	summary, err := a.dashboardService.GetDataSummary(uc.Uid)
	if err != nil {
		slog.Error("get summary error", err.Error())
	}

	dto.NewJsonResp(ctx).Success(summary)
}

// GetDashboardList 查询图表数据
func (a *DashboardApi) GetDashboardChart(ctx *gin.Context) {
	var req model.ReqDashboardChart
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	uc, err := GetUserClaim(ctx)
	if err != nil {
		slog.Error("get user claim err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.UnauthorizedError)
		return
	}

	data, err := a.dashboardService.GetChartData(req, uc.Uid)
	if err != nil {
		slog.Error("get chart data error", err.Error())
		if errors.Is(err, service.ErrUnSupportChartKey) {
			dto.NewJsonResp(ctx).Fail(dto.DashboardChartFailed, err)
			return
		}
		dto.NewJsonResp(ctx).Fail(dto.DashboardChartFailed)
		return
	}

	dto.NewJsonResp(ctx).Success(data)
}

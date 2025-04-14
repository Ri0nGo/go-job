package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/internal/model"
	"go-job/master/service"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type NodeApi struct {
	NodeService service.INodeService
}

func NewNodeApi(nodeService service.INodeService) *NodeApi {
	return &NodeApi{
		NodeService: nodeService,
	}
}

// RegisterRoutes 注册节点模块路由
func (a *NodeApi) RegisterRoutes(group *gin.RouterGroup) {
	nodeGroup := group.Group("/nodes")
	{
		nodeGroup.GET("", a.GetNodeList)
		nodeGroup.GET("/:id", a.GetNode)
		nodeGroup.POST("/add", a.AddNode)
		nodeGroup.PUT("/update", a.UpdateNode)
		nodeGroup.DELETE("/:id", a.DeleteNode)
	}
}

// GetNode 查询节点
func (a *NodeApi) GetNode(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	node, err := a.NodeService.GetNode(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dto.NewJsonResp(ctx).Success()
		return
	}
	if err != nil {
		slog.Error("get node err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(node)
}

// GetNodeList 查询节点列表
func (a *NodeApi) GetNodeList(ctx *gin.Context) {
	var page model.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	list, err := a.NodeService.GetNodeList(page)
	if err != nil {
		slog.Error("get node list err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeGetFailed)
		return
	}
	dto.NewJsonResp(ctx).Success(list)
}

// AddNode 添加节点
func (a *NodeApi) AddNode(ctx *gin.Context) {
	var req model.Node
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("add node params err", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	err := a.NodeService.AddNode(req)
	if err != nil {
		slog.Error("add node err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeAddFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// DeleteNode 删除节点
func (a *NodeApi) DeleteNode(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	if err := a.NodeService.DeleteNode(id); err != nil {
		slog.Error("delete node err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeDeleteFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

// UpdateNode 更新节点
func (a *NodeApi) UpdateNode(ctx *gin.Context) {
	var req model.Node
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}
	if _, err := a.NodeService.GetNode(req.Id); err != nil {
		dto.NewJsonResp(ctx).Fail(dto.NodeNotExist)
		return
	}
	if err := a.NodeService.UpdateNode(req); err != nil {
		slog.Error("update node err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeUpdateFailed)
		return
	}
	dto.NewJsonResp(ctx).Success()
}

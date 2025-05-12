package api

import (
	"github.com/gin-gonic/gin"
	"go-job/internal/dto"
	"go-job/node/service"
	"log/slog"
)

type NodeHandler struct {
	NodeService service.INodeService
}

func NewNodeHandler(NodService service.INodeService) *NodeHandler {
	return &NodeHandler{
		NodeService: NodService,
	}
}

// RegisterRoutes 注册job相关的路由, 遵循restful 风格
func (h *NodeHandler) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/install_ref", h.InstallRef)
	server.GET("info", h.NodeInfo)
}

// InstallRef 安装依赖
func (h *NodeHandler) InstallRef(ctx *gin.Context) {
	var req dto.ReqNodeRef
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("add Node bind json err:", "err", err)
		dto.NewJsonResp(ctx).Fail(dto.ParamsError)
		return
	}

	version, err := h.NodeService.InstallRef(ctx.Request.Context(), req)
	if err != nil {
		slog.Error("add Node error", "req", req, "err", err)
		dto.NewJsonResp(ctx).Fail(dto.NodeInstallRefFailed)
		return
	}

	dto.NewJsonResp(ctx).Success(map[string]string{
		"version":  version,
		"pkg_name": req.PkgName,
	})
}

func (h *NodeHandler) NodeInfo(ctx *gin.Context) {
	data := h.NodeService.GetNodeInfo(ctx.Request.Context())
	dto.NewJsonResp(ctx).Success(data)
}

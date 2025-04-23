package service

import "errors"

var (
	ErrCronExprParse      = errors.New("表达式无效")
	ErrSyncJobToNode      = errors.New("同步任务到节点失败")
	ErrSyncExecFileToNode = errors.New("同步执行文件到节点失败")
	ErrJobExtNotSupport   = errors.New("不支持的文件后缀")
	ErrNodeNotExists      = errors.New("节点不存在")
	ErrFileTooLarge       = errors.New("文件太大了")
	ErrInvalidAddress     = errors.New("填写的地址格式不合法，格式：Ip:Port")
	ErrJobUseCurrentNode  = errors.New("有任务依赖该节点，无法删除")
	ErrUserNotPermission  = errors.New("您的权限不足，暂无法使用此功能")
)

var returnErrList = []error{
	ErrCronExprParse,
	ErrSyncJobToNode,
	ErrSyncExecFileToNode,
	ErrJobExtNotSupport,
	ErrNodeNotExists,
	ErrFileTooLarge,
	ErrInvalidAddress,
	ErrJobUseCurrentNode,
	ErrUserNotPermission,
}

func IsRespErr(err error) bool {
	for _, knownErr := range returnErrList {
		if errors.Is(err, knownErr) {
			return true
		}
	}
	return false
}

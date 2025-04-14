package service

import "errors"

var (
	ErrCronExprParse      = errors.New("表达式无效")
	ErrSyncJobToNode      = errors.New("同步任务到节点失败")
	ErrSyncExecFileToNode = errors.New("同步执行文件到节点失败")
	ErrJobExtNotSupport   = errors.New("不支持的文件后缀")
	ErrNodeNotExists      = errors.New("节点不存在")
	ErrFileTooLarge       = errors.New("文件太大了")
)

var returnErrList = []error{
	ErrCronExprParse,
	ErrSyncJobToNode,
	ErrSyncExecFileToNode,
	ErrJobExtNotSupport,
	ErrNodeNotExists,
}

func IsRespErr(err error) bool {
	for _, knownErr := range returnErrList {
		if errors.Is(err, knownErr) {
			return true
		}
	}
	return false
}

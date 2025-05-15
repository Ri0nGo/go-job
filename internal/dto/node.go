package dto

import "go-job/internal/model"

type ReqNodeRef struct {
	Id      int                      `json:"id"`
	Type    model.NodeInstallRefType `json:"type"`
	PkgName string                   `json:"pkg_name"`
}

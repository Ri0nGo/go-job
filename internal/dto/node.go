package dto

import "go-job/internal/model"

type ReqNodeRef struct {
	Type    model.NodeInstallRefType `json:"type"`
	PkgName string                   `json:"pkg_name"`
	Version string                   `json:"version"`
}

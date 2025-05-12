package service

import (
	"context"
	"errors"
	"fmt"
	"go-job/internal/dto"
	"go-job/internal/model"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultTimeout = time.Second * 10

type INodeService interface {
	InstallRef(ctx context.Context, req dto.ReqNodeRef) (string, error)
	GetNodeInfo(ctx context.Context) map[string]any
}

type installRefInfo struct {
	pkgName string
	version string
}

type NodeService struct {
	timeout            time.Duration
	installRefHandlers map[model.NodeInstallRefType]func(ctx context.Context, info installRefInfo) (string, error)
}

func NewNodeService() *NodeService {
	s := &NodeService{
		timeout: defaultTimeout,
		// 返回对应的版本和错误
		installRefHandlers: make(map[model.NodeInstallRefType]func(ctx context.Context, info installRefInfo) (string, error)),
	}

	s.installRefHandlers[model.NodeInstallPyRefType] = s.installPyRef
	return s
}

func (s *NodeService) InstallRef(ctx context.Context, req dto.ReqNodeRef) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	fn, ok := s.installRefHandlers[req.Type]
	if !ok {
		return "", errors.New("invalid ref type")
	}
	return fn(ctx, installRefInfo{
		pkgName: req.PkgName,
		version: req.Version,
	})
}

func (s *NodeService) installPyRef(ctx context.Context, info installRefInfo) (string, error) {
	name := info.pkgName
	if len(info.version) > 0 {
		name = fmt.Sprintf("%s==%s", info.pkgName, info.version)
	}
	cmd := exec.CommandContext(ctx, "pip", "install", name)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("install py pkg failed: %w", err)
	}
	return getInstalledPyVersion(info.pkgName)
}

// getInstalledPyVersion 查询python包的版本
func getInstalledPyVersion(pkgName string) (string, error) {
	queryVersion := fmt.Sprintf("from importlib.metadata import version; print(version('%s'))", pkgName)
	cmd := exec.Command("python", "-c", queryVersion)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *NodeService) GetNodeInfo(ctx context.Context) map[string]any {
	return getPyInfo()
}

func getPyInfo() map[string]any {
	info := map[string]any{
		"version":  "unknown",
		"packages": "unknown",
	}

	// 获取 Python 版本
	versionCmd := exec.Command("python", "--version")
	versionOut, err := versionCmd.CombinedOutput()
	if err == nil {
		info["version"] = strings.TrimSpace(string(versionOut))
	}

	// 获取已安装包和版本
	listCmd := exec.Command("pip", "list", "--format=freeze")
	listOut, err := listCmd.Output()
	if err != nil {
		return info
	}

	var pkgs []string
	for _, s := range strings.Split(string(listOut), "\n") {
		pkgs = append(pkgs, strings.Replace(s, "\r", "", -1))
	}
	info["packages"] = pkgs
	return info
}

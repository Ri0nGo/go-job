package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-job/internal/dto"
	"go-job/internal/model"
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

	s.installRefHandlers[model.NodeInstallPyRefType] = installPyRef
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
	})
}

func installPyRef(ctx context.Context, info installRefInfo) (string, error) {
	name := info.pkgName
	cmd := exec.CommandContext(ctx, "pip", "install", name)

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("install py pkg failed: %w, stderr: %s, stdout: %s",
			err, stderr.String(), stdout.String())
	}
	return getInstalledPyVersion(info.pkgName)
}

// getInstalledPyVersion 查询python包的版本
func getInstalledPyVersion(pkgName string) (string, error) {
	script := fmt.Sprintf(`from importlib.metadata import version; print(version("%s"))`, pkgName)
	cmd := exec.Command("python", "-c", script)

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("get py pkg failed: %w, stderr: %s, stdout: %s",
			err, stderr.String(), stdout.String())
	}
	return strings.TrimSpace(stdout.String()), nil
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

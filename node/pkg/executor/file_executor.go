package executor

import (
	"errors"
	"fmt"
	"go-job/internal/models"
	"go-job/node/pkg/config"
	"os/exec"
	"path/filepath"
)

type FileExecutor struct {
	ext            string
	fileName       string
	OnStatusChange func(status models.JobStatus) // 注册回调事件， 后续可以优化为channel的方式接收结果
}

func (f *FileExecutor) Run() {
	f.OnStatusChange(models.Running)
	if err := f.Execute(); err != nil {
		f.OnStatusChange(models.Failed)
		return
	}
	f.OnStatusChange(models.Success)

}

func (f *FileExecutor) Execute() error {
	var (
		output []byte
		err    error
	)
	switch f.ext {
	case ".py":
		output, err = f.execFile()
	default:
		output = []byte("不支持的文件类型")
		err = errors.New("不支持的文件类型")
	}

	// todo 回传结果和执行状态
	fmt.Println("exec result: ", string(output), err)
	return err
}

func NewFileExecutor(fileName string) *FileExecutor {
	return &FileExecutor{
		ext:      filepath.Ext(fileName),
		fileName: fileName,
	}
}

func (f *FileExecutor) SetOnStatusChange(fn func(status models.JobStatus)) {
	f.OnStatusChange = fn
}

func (f *FileExecutor) execFile() (output []byte, err error) {
	// 执行文件，这是一次性捕获所有输出，无法实现实时捕获，
	execFilePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir, f.fileName)
	// TODO 后续需要修改为实时捕获
	cmd := exec.Command("python3", execFilePath)
	return cmd.CombinedOutput()
}

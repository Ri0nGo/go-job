package executor

import (
	"errors"
	"fmt"
	"go-job/internal/model"
	"go-job/node/pkg/config"
	"os/exec"
	"path/filepath"
)

type FileExecutor struct {
	id             int
	name           string
	ext            string
	fileName       string
	OnStatusChange func(status model.JobStatus) // 注册回调事件， 后续可以优化为channel的方式接收结果
}

func (f *FileExecutor) Run() {
	f.OnStatusChange(model.Running)
	if err := f.Execute(); err != nil {
		f.OnStatusChange(model.Failed)
		return
	}
	f.OnStatusChange(model.Success)

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
	fmt.Printf("id: %d, name: %s, exec result: %s, err: %v\n",
		f.id, f.name, string(output), err)

	// todo 待实现
	go f.sendResultToMaster()
	return err
}

func NewFileExecutor(id int, name, fileName string) *FileExecutor {
	return &FileExecutor{
		id:       id,
		name:     name,
		ext:      filepath.Ext(fileName),
		fileName: fileName,
	}
}

func (f *FileExecutor) SetOnStatusChange(fn func(status model.JobStatus)) {
	f.OnStatusChange = fn
}

func (f *FileExecutor) execFile() (output []byte, err error) {
	// 执行文件，这是一次性捕获所有输出，无法实现实时捕获，
	execFilePath := fmt.Sprintf("%s/%s", config.App.Data.UploadJobDir, f.fileName)
	// TODO 后续需要修改为实时捕获
	cmd := exec.Command("python3", execFilePath)
	return cmd.CombinedOutput()
}

// sendResultToMaster 回传结果到master
func (f *FileExecutor) sendResultToMaster() {

}

package executor

import (
	"errors"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/node/pkg/config"
	"os/exec"
	"path/filepath"
	"time"
)

var defaultOutputLen = 1024

type FileExecutor struct {
	id             int
	name           string
	ext            string
	fileName       string
	runningStatus  model.JobStatus
	startExecTime  time.Time
	endExecTime    time.Time
	onResultChange func(result model.JobExecResult) // 注册回调事件， 后续可以优化为channel的方式接收结果
}

func (f *FileExecutor) Run() {
	f.BeforeExecute()
	output, err := f.Execute()
	f.AfterExecute(err)
	if f.onResultChange != nil {
		f.onResultChange(f.buildJobExecResult(output, err))
	}

}

func (f *FileExecutor) BeforeExecute() {
	f.startExecTime = time.Now()
	f.runningStatus = model.Running
}

func (f *FileExecutor) AfterExecute(err error) {
	if err != nil {
		f.runningStatus = model.Failed
	} else {
		f.runningStatus = model.Success
	}
	f.endExecTime = time.Now()
}

func (f *FileExecutor) buildJobExecResult(output string, err error) model.JobExecResult {
	runes := []rune(output)
	if len(output) > defaultOutputLen {
		output = string(runes[:defaultOutputLen-3]) + "..."
	}
	result := model.JobExecResult{
		StartTime: f.startExecTime.Unix(),
		EndTime:   f.endExecTime.Unix(),
		Duration:  f.endExecTime.Sub(f.startExecTime).Seconds(),
		Status:    f.runningStatus,
		Output:    output,
		Error:     utils.ErrorToString(err),
	}
	return result
}

func (f *FileExecutor) Execute() (string, error) {
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

	//fmt.Printf("id: %d, name: %s, exec result: %s, err: %v\n",
	//	f.id, f.name, string(output), err)

	return string(output), err
}

func (f *FileExecutor) OnResultChange(fn func(result model.JobExecResult)) {
	f.onResultChange = fn
}

func NewFileExecutor(id int, name, fileName string) *FileExecutor {
	return &FileExecutor{
		id:       id,
		name:     name,
		ext:      filepath.Ext(fileName),
		fileName: fileName,
	}
}

func (f *FileExecutor) execFile() (output []byte, err error) {
	// 执行文件，这是一次性捕获所有输出，无法实现实时捕获，
	execFilePath := filepath.Join(config.App.Data.UploadJobDir, f.fileName)
	// TODO 后续需要修改为实时捕获
	cmd := exec.Command("python", execFilePath)
	return cmd.CombinedOutput()
}

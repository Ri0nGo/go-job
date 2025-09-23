package executor

import (
	"errors"
	"fmt"
	"go-job/internal/model"
	"testing"
)

type MockFileHandler struct {
}

func (m MockFileHandler) Run() {
	m.Execute()
}

func (m MockFileHandler) Execute() (string, error) {
	return "", errors.New("mock a error")
}

func (m MockFileHandler) OnResultChange(f func(result model.JobExecResult)) {
	fmt.Println("OnResultChange")
}

func (m MockFileHandler) ResultCallback(output string, err error) {
	fmt.Println("ResultCallback")
}

func (m MockFileHandler) AfterExecute(err error) {
	fmt.Println("exec after")
}

func (m MockFileHandler) BeforeExecute() {
	fmt.Println("exec before")
}

func TestRetryExecutor_Run(t *testing.T) {
	mockH := MockFileHandler{}
	executor := NewRetryExecutor(mockH, 3)
	executor.Run()
}

package executor

import "fmt"

type FileExecutor struct {
}

func (f *FileExecutor) Run() {
	fmt.Printf("Running file executor\n")
}

func NewFileExecutor() *FileExecutor {
	return &FileExecutor{}
}

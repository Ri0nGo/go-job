package utils

import (
	"fmt"
	"os"
)

// EnsureDir 目录不存在则创建
func EnsureDir(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s 已存在但不是一个目录", path)
	}
	return nil
}

func ErrorToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

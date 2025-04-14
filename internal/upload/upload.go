package upload

import (
	"errors"
	"github.com/Ri0nGo/gokit/slice"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrFileTooLarge        = errors.New("file too large")
	ErrFileExtNotSupported = errors.New("file extension not supported")
)

// FileMeta 文件信息
type FileMeta struct {
	Filename     string    `json:"filename"`       // 原始文件名
	UUIDFileName string    `json:"uuid_file_name"` // 修改后的文件名
	Size         int       `json:"size"`
	UploadTime   time.Time `json:"upload_time"`
}

// ValidatorOptions 文件信息校验器
type ValidatorOptions func(fileMeta FileMeta) error

// Option 修改defaultFileUpload的直接
type Option func(*FileUpload)

type FileUpload struct {
	mux   sync.RWMutex
	files map[string]FileMeta
	exts  []string
	size  int
}

var defaultFu = &FileUpload{
	files: make(map[string]FileMeta),
	exts:  []string{".py"},
	size:  5 * 1024 * 1024,
}

// ----------- FileUpload ----------- //

// SetFileMeta 存储一个FileMeta
func SetFileMeta(uuid string, file FileMeta) {
	defaultFu.mux.Lock()
	defer defaultFu.mux.Unlock()
	defaultFu.files[uuid] = file
}

// GetFileMeta 获取一个FileMeta
func GetFileMeta(uuid string) (FileMeta, bool) {
	defaultFu.mux.RLock()
	defer defaultFu.mux.RUnlock()
	file, ok := defaultFu.files[uuid]
	return file, ok
}

func DeleteFileMeta(uuid string) {
	defaultFu.mux.Lock()
	defer defaultFu.mux.Unlock()
	delete(defaultFu.files, uuid)
}

func ExtsOpt(exts []string) Option {
	return func(f *FileUpload) {
		f.exts = exts
	}
}

func SizeOpt(size int) Option {
	return func(f *FileUpload) {
		f.size = size
	}
}

func FileUploadOpts(opts ...Option) {
	for _, opt := range opts {
		opt(defaultFu)
	}
}

// ----------- FileMeta 相关 ----------- //

func FileExtValidator(fileMeta FileMeta) error {
	if ok := slice.Contains(defaultFu.exts, filepath.Ext(fileMeta.Filename)); !ok {
		return ErrFileExtNotSupported
	}
	return nil
}

func FileSizeValidator(fileMeta FileMeta) error {
	if fileMeta.Size > defaultFu.size {
		return ErrFileTooLarge
	}
	return nil
}

func ValidatorFileOpts(fileMeta FileMeta, opts ...ValidatorOptions) error {
	for _, opt := range opts {
		if err := opt(fileMeta); err != nil {
			return err
		}
	}
	return nil
}

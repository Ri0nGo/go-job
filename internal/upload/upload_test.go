package upload

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatorFileOpts(t *testing.T) {
	testCases := []struct {
		name      string
		input     FileMeta
		validator []ValidatorOptions
		err       error
	}{
		{
			name: "extension not supported",
			input: FileMeta{
				Filename: "test.txt",
			},
			validator: []ValidatorOptions{FileExtValidator},
			err:       ErrFileExtNotSupported,
		}, {
			name: "size to large",
			input: FileMeta{
				Size: 10 * 1024 * 1024,
			},
			validator: []ValidatorOptions{FileSizeValidator},
			err:       ErrFileTooLarge,
		}, {
			name: "ext success",
			input: FileMeta{
				Filename: "test.py",
			},
			validator: []ValidatorOptions{FileExtValidator},
			err:       nil,
		}, {
			name: "size success",
			input: FileMeta{
				Size: 1024 * 1024,
			},
			validator: []ValidatorOptions{FileSizeValidator},
			err:       nil,
		}, {
			name: "size & ext success",
			input: FileMeta{
				Filename: "test.py",
				Size:     1024 * 1024,
			},
			validator: []ValidatorOptions{FileSizeValidator, FileExtValidator},
			err:       nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatorFileOpts(tc.input, tc.validator...)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestFileUploadOpts(t *testing.T) {
	exts := []string{".py", ".java", ".php"}
	size := 1024 * 1024
	FileUploadOpts(
		ExtsOpt(exts),
		SizeOpt(size))
	assert.Equal(t, defaultFu.exts, exts)
	assert.Equal(t, defaultFu.size, size)
}

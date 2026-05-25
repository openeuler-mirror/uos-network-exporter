package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"uos_network_exporter/pkg/utils"
)

var (
	defaultMaxFiles = 5
)

type FileRotator struct {
	basePath  string
	maxSize   int64
	maxAge    time.Duration
	current   *os.File
	size      int64
	startTime time.Time
	keepFiles int
}

func NewFileRotator(basePath string, maxSize int64, maxAge time.Duration) *FileRotator {
	dir := filepath.Dir(basePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0750); err != nil {
			fmt.Printf("Warning: Failed to create the log directory: %v\n", err)
		}
	}
	return &FileRotator{
		basePath:  basePath,
		maxSize:   maxSize,
		maxAge:    maxAge,
		keepFiles: defaultMaxFiles,
	}
}

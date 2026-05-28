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

func (fr *FileRotator) Write(p []byte) (n int, err error) {
	err = fr.setupCurrent()
	if err != nil {
		return 0, err
	}
	if fr.shouldRotate() {
		err = fr.rotate()
		if err != nil {
			return 0, err
		}
	}
	n, err = fr.current.Write(p)
	if err != nil {
		return n, err
	}
	fr.size += int64(n)
	return n, nil
}

func (fr *FileRotator) setupCurrent() error {
	if fr.current == nil {
		fileinfo, err := os.Stat(fr.basePath)
		if err == nil {
			fr.current, err = os.OpenFile(fr.basePath, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			fr.size = fileinfo.Size()
			fr.startTime = fileinfo.ModTime()
		} else if os.IsNotExist(err) {
			fr.current, err = os.Create(fr.basePath)
			if err != nil {
				return err
			}
			fr.startTime = time.Now()
		} else {
			return err
		}
	}
	return nil
}

func (fr *FileRotator) Close() error {
	if fr.current != nil {
		return fr.current.Close()
	}
	return nil
}

func (fr *FileRotator) shouldRotate() bool {
	if fr.size > fr.maxSize || time.Now().Sub(fr.startTime) > fr.maxAge {
		return true
	}
	return false
}

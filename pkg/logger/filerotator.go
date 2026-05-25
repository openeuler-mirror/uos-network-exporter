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

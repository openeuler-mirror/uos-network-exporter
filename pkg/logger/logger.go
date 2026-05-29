package logger

import (
	"strings"
	"time"

	formatter "gitee.com/weidongkl/logrus-formatter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Level   string        `yaml:"level"`
	LogPath string        `yaml:"log_path"`
	MaxSize string        `yaml:"max_size"`
	MaxAge  time.Duration `yaml:"max_age"`
}

type fileLogConfig struct {
	FileRotator *FileRotator
	level       string
}

func NewConfig(level, logPath string, maxSize int64, maxAge time.Duration) fileLogConfig {
	return fileLogConfig{
		level:       level,
		FileRotator: NewFileRotator(logPath, maxSize, maxAge),
	}
}

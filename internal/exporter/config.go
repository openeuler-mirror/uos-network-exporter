package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	Configfile        *string
	NetworkConfigfile *string
	DefaultConfig     = Config{
		Logging: logger.Config{
			Level:   "debug",
			LogPath: "./network_exporter.log",
			MaxSize: "10MB",
			MaxAge:  time.Hour * 24 * 7},
		Address:     "127.0.0.1",
		Port:        9118,
		MetricsPath: "/metrics",
	}
)

func init() {
	Configfile = kingpin.Flag("config", "Configuration file").
		Short('c').
		Default("/etc/uos-exporter/network-exporter.yaml").
		String()
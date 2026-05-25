package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	ScrapeUrl       *string
	Insecure        *bool
	DefaultSettings = Settings{
		Insecure: false,
	}
)

func init() {
	ScrapeUrl = kingpin.Flag("scrape_uri", "Scrape URI").Short('s').String()
	Insecure = kingpin.Flag("insecure", "Ignore server certificate").Bool()
}

type Settings struct {
	ScrapeUri string `yaml:"scrape_uri"`

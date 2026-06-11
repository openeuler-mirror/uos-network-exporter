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
	Insecure  bool   `yaml:"insecure"`
}

type Targets []struct {
	Name     string   `yaml:"name" json:"name"`
	Host     string   `yaml:"host" json:"host"`
	Port     string   `yaml:"port" json:"port"`
	Type     string   `yaml:"type" json:"type"`
	Proxy    string   `yaml:"proxy" json:"proxy"`
	Probe    []string `yaml:"probe" json:"probe"`
	SourceIp string   `yaml:"source_ip" json:"source_ip"`
	Labels   extraKV  `yaml:"labels,omitempty" json:"labels,omitempty"`
}

type HTTPGet struct {
	Interval duration `yaml:"interval" json:"interval"`
	Timeout  duration `yaml:"timeout" json:"timeout"`
}

type TCP struct {
	Interval duration `yaml:"interval" json:"interval"`
	Timeout  duration `yaml:"timeout" json:"timeout"`
}

type MTR struct {
	Interval duration `yaml:"interval" json:"interval"`
	Timeout  duration `yaml:"timeout" json:"timeout"`
	MaxHops  int      `yaml:"max-hops" json:"max-hops"`
	Count    int      `yaml:"count" json:"count"`
}

type ICMP struct {
	Interval duration `yaml:"interval" json:"interval"`
	Timeout  duration `yaml:"timeout" json:"timeout"`
	Count    int      `yaml:"count" json:"count"`
}

type Conf struct {
	Refresh           duration `yaml:"refresh" json:"refresh"`
	Nameserver        string   `yaml:"nameserver" json:"nameserver"`

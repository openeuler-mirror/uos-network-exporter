package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var defaultReg *Registry

func init() {
	defaultReg = NewRegistry()
}

type Registry struct {
	metrics []Metric
	mu      sync.RWMutex
}

func Register(metric Metric) {
	defaultReg.Register(metric)
}

func RegisterPrometheus(reg *prometheus.Registry) {
	reg.MustRegister(defaultReg)
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: []Metric{},
	}
}
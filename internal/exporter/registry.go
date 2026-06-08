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
package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// metricInfo 存储单个metric的信息
type metricInfo struct {
	desc       *prometheus.Desc
	labelNames []string
	values     map[string]float64
}
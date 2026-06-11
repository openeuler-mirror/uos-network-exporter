package metrics

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/mtr"
	"github.com/prometheus/client_golang/prometheus"
)

type MTRCacheEntry struct {
	result    *mtr.MtrResult
	timestamp time.Time
}

type MTRMetrics struct {
	*baseMetrics
	logger   *slog.Logger
	resolver *net.Resolver
	cache    map[string]*MTRCacheEntry
	cacheMux sync.RWMutex
	cacheTTL time.Duration
}

func NewMTRMetrics(logger *slog.Logger, resolver *net.Resolver) *MTRMetrics {
	base := newBaseMetrics("mtr")
	base.addMetric("rtt_seconds", "Round Trip Time in seconds", []string{"name", "target", "ttl", "path", "type"})
	base.addMetric("rtt_snt_count", "Round Trip Send Package Total", []string{"name", "target", "ttl", "path"})
	base.addMetric("rtt_snt_fail_count", "Round Trip Send Package Fail Total", []string{"name", "target", "ttl", "path"})
	base.addMetric("rtt_snt_seconds", "Round Trip Send Package Time Total", []string{"name", "target", "ttl", "path"})
	base.addMetric("hops", "Number of route hops", []string{"name", "target"})
	base.addMetric("targets", "Number of active targets", nil)
	base.addMetric("up", "Exporter state", nil)
	return &MTRMetrics{
		baseMetrics: base,
		logger:      logger,
		resolver:    resolver,
		cache:       make(map[string]*MTRCacheEntry),
		cacheTTL:    30 * time.Second,
	}
}

func (m *MTRMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.baseMetrics.Describe(ch)
}

func (m *MTRMetrics) CollectMetrics(ch chan<- prometheus.Metric) {
	m.baseMetrics.Collect(ch)
}

func (m *MTRMetrics) getCachedResult(key string) *mtr.MtrResult {
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()
	if entry, exists := m.cache[key]; exists {
		if time.Since(entry.timestamp) < m.cacheTTL {
			return entry.result
		}

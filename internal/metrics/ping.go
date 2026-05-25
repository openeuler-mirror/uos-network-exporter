package metrics

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/ping"

	"github.com/prometheus/client_golang/prometheus"
)

// PingCacheEntry PING缓存条目
type PingCacheEntry struct {
	result    *ping.PingResult
	timestamp time.Time
}

// PingMetrics ping相关的metrics
type PingMetrics struct {
	*baseMetrics
	logger   *slog.Logger
	resolver *net.Resolver
	icmpID   *common.IcmpID
	cache    map[string]*PingCacheEntry
	cacheMux sync.RWMutex
	cacheTTL time.Duration
}
// NewPingMetrics 创建新的ping metrics实例
func NewPingMetrics(logger *slog.Logger, resolver *net.Resolver) *PingMetrics {
	base := newBaseMetrics("ping")
	base.addMetric("status", "Ping Status", []string{"name", "target", "target_ip"})
	base.addMetric("rtt_seconds", "Round Trip Time", []string{"name", "target", "target_ip", "type"})

	base.addMetric("rtt_snt_count", "Packet sent count", []string{"name", "target", "target_ip"})
	base.addMetric("rtt_snt_fail_count", "Packet sent fail count", []string{"name", "target", "target_ip"})
	base.addMetric("rtt_snt_seconds", "Packet sent time total", []string{"name", "target", "target_ip"})
	base.addMetric("loss_percent", "Packet loss in percent", []string{"name", "target", "target_ip"})
	base.addMetric("targets", "Number of active targets", nil)
	base.addMetric("up", "Exporter state", nil)

	return &PingMetrics{
		baseMetrics: base,
		logger:      logger,
		resolver:    resolver,
		icmpID:      &common.IcmpID{},
		cache:       make(map[string]*PingCacheEntry),
		cacheTTL:    15 * time.Second,
	}
}

func (p *PingMetrics) Describe(ch chan<- *prometheus.Desc) {
	p.baseMetrics.Describe(ch)
}

func (p *PingMetrics) CollectMetrics(ch chan<- prometheus.Metric) {
	p.baseMetrics.Collect(ch)
}

func (p *PingMetrics) getCachedResult(key string) *ping.PingResult {
	p.cacheMux.RLock()
	defer p.cacheMux.RUnlock()
	if entry, exists := p.cache[key]; exists {
		if time.Since(entry.timestamp) < p.cacheTTL {
			return entry.result

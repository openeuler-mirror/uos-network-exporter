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

package exporter

import (
	"context"
	"log/slog"
	"net"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

// NetworkExporter 网络监控导出器
type NetworkExporter struct {
	logger      *slog.Logger
	config      *config.SafeConfig
	resolver    *net.Resolver
	pingMetrics *metrics.PingMetrics
	tcpMetrics  *metrics.TCPMetrics
	httpMetrics *metrics.HTTPMetrics
	mtrMetrics  *metrics.MTRMetrics
}

// NewNetworkExporter 创建新的网络导出器
func NewNetworkExporter(logger *slog.Logger, cfg *config.SafeConfig) *NetworkExporter {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, address)
		},
	}

	return &NetworkExporter{
		logger:   logger,
		config:   cfg,
		resolver: resolver,
		pingMetrics: metrics.NewPingMetrics(logger, resolver),
		tcpMetrics:  metrics.NewTCPMetrics(logger, resolver),
		httpMetrics: metrics.NewHTTPMetrics(logger, resolver),
		mtrMetrics:  metrics.NewMTRMetrics(logger, resolver),
	}
}

// Describe 实现prometheus.Collector接口
func (ne *NetworkExporter) Describe(ch chan<- *prometheus.Desc) {
	ne.pingMetrics.Describe(ch)
	ne.tcpMetrics.Describe(ch)
	ne.httpMetrics.Describe(ch)
	ne.mtrMetrics.Describe(ch)
}
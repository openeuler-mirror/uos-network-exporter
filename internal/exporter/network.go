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
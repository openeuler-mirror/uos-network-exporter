package metrics

import (
	"context"
	"log/slog"
	"net"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/tcp"
	"github.com/prometheus/client_golang/prometheus"
)

type TCPMetrics struct {
	*baseMetrics
	logger   *slog.Logger
	resolver *net.Resolver
}

func NewTCPMetrics(logger *slog.Logger, resolver *net.Resolver) *TCPMetrics {
	base := newBaseMetrics("tcp")
	base.addMetric("connection_seconds", "Connection time in seconds", []string{"name", "target", "target_ip", "source_ip", "port"})
	base.addMetric("connection_status", "Connection Status", []string{"name", "target", "target_ip", "source_ip", "port"})
	base.addMetric("targets", "Number of active targets", nil)
	base.addMetric("up", "Exporter state", nil)
	return &TCPMetrics{
		baseMetrics: base,
		logger:      logger,
		resolver:    resolver,
	}
}

func (t *TCPMetrics) Describe(ch chan<- *prometheus.Desc) {
	t.baseMetrics.Describe(ch)
}

func (t *TCPMetrics) CollectMetrics(ch chan<- prometheus.Metric) {
	t.baseMetrics.Collect(ch)
}

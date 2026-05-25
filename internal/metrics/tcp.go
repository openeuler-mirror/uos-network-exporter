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

func (t *TCPMetrics) Collect(cfg *config.NetworkConfig) {
	if cfg == nil {
		t.logger.Warn("TCP config is nil")
		return
	}
	t.setMetric("up", 1)
	targetCount := 0
	for _, target := range cfg.Targets {
		if target.Type != "TCP" {
			continue
		}
		targetCount++
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		ipAddrs, err := common.DestAddrs(ctx, target.Host, t.resolver, 5*time.Second)
		cancel()
		if err != nil || len(ipAddrs) == 0 {
			t.logger.Warn("Failed to resolve target", "host", target.Host, "error", err)
			t.setMetricWithLabels("connection_status", 0, map[string]string{
				"name": target.Name, "target": target.Host, "target_ip": "unknown",
				"source_ip": target.SourceIp, "port": target.Port,
			})
			continue
		}
		for _, ip := range ipAddrs {
			result := tcp.TestTCPPort(target.Host, ip, target.Port, target.SourceIp, cfg.TCP.Timeout.Duration())
			labels := map[string]string{
				"name": target.Name, "target": target.Host, "target_ip": ip,
				"source_ip": result.SrcIp, "port": target.Port,
			}
			if result.Success {
				t.setMetricWithLabels("connection_status", 1, labels)

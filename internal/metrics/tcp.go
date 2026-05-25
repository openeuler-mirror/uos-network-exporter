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

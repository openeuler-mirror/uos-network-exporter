package metrics

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"uos_network_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	*baseMetrics
	logger   *slog.Logger
	resolver *net.Resolver
}

func NewHTTPMetrics(logger *slog.Logger, resolver *net.Resolver) *HTTPMetrics {
	base := newBaseMetrics("http")
	base.addMetric("get_seconds", "HTTP Get Drill Down time in seconds", []string{"name", "target", "type"})
	base.addMetric("get_content_bytes", "HTTP Get Content Size in bytes", []string{"name", "target"})
	base.addMetric("get_status", "HTTP Get Status", []string{"name", "target"})
	base.addMetric("get_targets", "Number of active targets", nil)
	base.addMetric("get_up", "Exporter state", nil)
	return &HTTPMetrics{
		baseMetrics: base,
		logger:      logger,
		resolver:    resolver,
	}
}

func (h *HTTPMetrics) Describe(ch chan<- *prometheus.Desc) {
	h.baseMetrics.Describe(ch)
}

func (h *HTTPMetrics) CollectMetrics(ch chan<- prometheus.Metric) {
	h.baseMetrics.Collect(ch)
}


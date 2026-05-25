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

func (h *HTTPMetrics) Collect(cfg *config.NetworkConfig) {
	if cfg == nil {
		h.logger.Warn("HTTPGet config is nil")
		return
	}
	h.setMetric("get_up", 1)
	targetCount := 0
	for _, target := range cfg.Targets {
		if target.Type != "HTTPGet" {
			continue
		}
		targetCount++
		client := &http.Client{
			Timeout: time.Second * 3,
		}
		start := time.Now()
		resp, err := client.Get(target.Host)
		duration := time.Since(start)
		success := 0.0
		statusCode := 0.0
		contentLength := 0.0
		if err == nil {
			success = 1.0
			statusCode = float64(resp.StatusCode)
			contentLength = float64(resp.ContentLength)
			if err := resp.Body.Close(); err != nil {
				h.logger.Warn("failed to close response body", "error", err)
			}
		}
		baseLabels := map[string]string{
			"name": target.Name, "target": target.Host,
		}
		h.setMetricWithLabels("get_status", statusCode, baseLabels)
		h.setMetricWithLabels("get_content_bytes", contentLength, baseLabels)
		h.setMetricWithLabels("get_seconds", duration.Seconds(), map[string]string{
			"name": target.Name, "target": target.Host, "type": "Total",
		})
		h.setMetric("get_up", success)
	}
	h.setMetric("get_targets", float64(targetCount))
}

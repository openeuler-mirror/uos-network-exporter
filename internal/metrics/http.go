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

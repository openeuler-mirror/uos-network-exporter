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
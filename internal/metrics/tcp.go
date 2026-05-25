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

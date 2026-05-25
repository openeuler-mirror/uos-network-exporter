package metrics

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/ping"

	"github.com/prometheus/client_golang/prometheus"
)

// PingCacheEntry PING缓存条目
type PingCacheEntry struct {
	result    *ping.PingResult
	timestamp time.Time
}

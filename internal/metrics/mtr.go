package metrics

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"uos_network_exporter/config"
	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/mtr"
	"github.com/prometheus/client_golang/prometheus"
)

type MTRCacheEntry struct {
	result    *mtr.MtrResult
	timestamp time.Time
}

type MTRMetrics struct {
	*baseMetrics
	logger   *slog.Logger
	resolver *net.Resolver
	cache    map[string]*MTRCacheEntry
	cacheMux sync.RWMutex
	cacheTTL time.Duration
}

func NewMTRMetrics(logger *slog.Logger, resolver *net.Resolver) *MTRMetrics {

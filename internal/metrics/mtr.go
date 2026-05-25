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

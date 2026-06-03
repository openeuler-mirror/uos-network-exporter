package mtr

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/icmp"
)

// ConcurrentMTROptions concurrent MTR configuration options
type ConcurrentMTROptions struct {
	MaxWorkers     int
	BatchSize      int
	EarlyStop      bool
	ProgressReport bool
	Timeout        time.Duration
}

// DefaultConcurrentMTROptions returns default concurrent MTR config
func DefaultConcurrentMTROptions() *ConcurrentMTROptions {
	return &ConcurrentMTROptions{
		MaxWorkers:     10,
		BatchSize:      5,
		EarlyStop:      true,
		ProgressReport: false,
		Timeout:        30 * time.Second,
	}
}

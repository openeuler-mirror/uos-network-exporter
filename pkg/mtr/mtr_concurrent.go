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

// RunMTRConcurrent execute high-performance concurrent MTR route tracing
func RunMTRConcurrent(destAddr, srcAddr string, timeout time.Duration, maxHops, count int, options *ConcurrentMTROptions) *MtrResult {
	if options == nil {
		options = DefaultConcurrentMTROptions()
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	result := &MtrResult{
		DestAddr:      destAddr,
		Hops:          []common.IcmpHop{},
		HopSummaryMap: make(map[string]*common.IcmpSummary),
	}

	icmpID := &common.IcmpID{}
	pid := int(icmpID.Get())

	mtrReturns := make([]*safeReturn, maxHops+1)
	for i := 1; i < maxHops+1; i++ {
		mtrReturns[i] = &safeReturn{
			mtrReturn: &MtrReturn{
				ttl:       i,
				host:      "unknown",
				succSum:   0,
				success:   false,
				lastTime:  time.Duration(0),
				sumTime:   time.Duration(0),
				bestTime:  time.Duration(0),
				worstTime: time.Duration(0),
				avgTime:   time.Duration(0),
				allTime:   make([]time.Duration, 0, count),
			},
		}
	}

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

	taskChan := make(chan task, options.BatchSize*options.MaxWorkers)
	go func() {
		defer close(taskChan)
		for snt := 0; snt < count; snt++ {
			for ttl := 1; ttl < maxHops; ttl++ {
				select {
				case taskChan <- task{ttl: ttl, snt: snt}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	var wg sync.WaitGroup
	var seqCounter int64
	reachedDestination := int32(0)

	workerCount := options.MaxWorkers
	if workerCount > maxHops {
		workerCount = maxHops
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case task, ok := <-taskChan:
					if !ok {
						return
					}

					if options.EarlyStop && atomic.LoadInt32(&reachedDestination) == 1 && task.ttl > int(atomic.LoadInt32(&reachedDestination)) {
						continue
					}

					processTask(ctx, task, destAddr, srcAddr, pid, timeout, &seqCounter, mtrReturns, &reachedDestination)

				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	wg.Wait()

	buildResult(result, mtrReturns, destAddr, count)

	return result
}

// safeReturn thread-safe return structure
type safeReturn struct {
	mu        sync.RWMutex
	mtrReturn *MtrReturn
}

// task represents one ICMP test task
type task struct {
	ttl int
	snt int
}

// processTask process a single task
func processTask(ctx context.Context, t task, destAddr, srcAddr string, pid int, timeout time.Duration, seqCounter *int64, mtrReturns []*safeReturn, reachedDestination *int32) {
	seq := int(atomic.AddInt64(seqCounter, 1))

	select {
	case <-ctx.Done():
		return
	default:
	}

	hopReturn, err := icmp.Icmp(destAddr, srcAddr, t.ttl, pid, timeout, seq, false)
	if err != nil || !hopReturn.Success {
		return
	}

	safeRet := mtrReturns[t.ttl]
	safeRet.mu.Lock()
	defer safeRet.mu.Unlock()

	ret := safeRet.mtrReturn
	ret.host = hopReturn.Addr
	ret.lastTime = hopReturn.Elapsed
	ret.allTime = append(ret.allTime, hopReturn.Elapsed)
	ret.succSum++

	if ret.worstTime == time.Duration(0) || hopReturn.Elapsed > ret.worstTime {
		ret.worstTime = hopReturn.Elapsed
	}
	if ret.bestTime == time.Duration(0) || hopReturn.Elapsed < ret.bestTime {
		ret.bestTime = hopReturn.Elapsed
	}
	ret.sumTime += hopReturn.Elapsed
	ret.avgTime = time.Duration((int64)(ret.sumTime/time.Microsecond)/(int64)(ret.succSum)) * time.Microsecond
	ret.success = true

	if common.IsEqualIP(hopReturn.Addr, destAddr) {
		atomic.StoreInt32(reachedDestination, safeIntToInt32(t.ttl))
	}
}

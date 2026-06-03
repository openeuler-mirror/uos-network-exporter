package mtr

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"uos_network_exporter/pkg/common"
	"uos_network_exporter/pkg/icmp"
)

// RunMTR execute real MTR route tracing (concurrent version)
func RunMTR(destAddr, srcAddr string, timeout time.Duration, maxHops, count int) *MtrResult {
	result := &MtrResult{
		DestAddr:      destAddr,
		Hops:          []common.IcmpHop{},
		HopSummaryMap: make(map[string]*common.IcmpSummary),
	}

	// use ICMP ID
	icmpID := &common.IcmpID{}
	pid := int(icmpID.Get())

	mtrReturns := make([]*MtrReturn, maxHops+1)
	var mtrMutex sync.Mutex

	var seqCounter int64

	type task struct {
		ttl int
		snt int
	}

	taskChan := make(chan task, maxHops*count)

	// generate all tasks
	for snt := 0; snt < count; snt++ {
		for ttl := 1; ttl < maxHops; ttl++ {
			taskChan <- task{ttl: ttl, snt: snt}
		}
	}
	close(taskChan)

	workerCount := maxHops
	if workerCount > 20 {
		workerCount = 20
	}

	var wg sync.WaitGroup
	hasPermissionError := false
	var permissionMutex sync.Mutex

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskChan {
				ttl := task.ttl
				seq := int(atomic.AddInt64(&seqCounter, 1))

				mtrMutex.Lock()
				if mtrReturns[ttl] == nil {
					mtrReturns[ttl] = &MtrReturn{
						ttl:       ttl,
						host:      "unknown",
						succSum:   0,
						success:   false,
						lastTime:  time.Duration(0),
						sumTime:   time.Duration(0),
						bestTime:  time.Duration(0),
						worstTime: time.Duration(0),
						avgTime:   time.Duration(0),
						allTime:   make([]time.Duration, 0),
					}
				}
				mtrMutex.Unlock()

				hopReturn, err := icmp.Icmp(destAddr, srcAddr, ttl, pid, timeout, seq, false)

				if err != nil {
					permissionMutex.Lock()
					if !hasPermissionError {
						hasPermissionError = true
					}
					permissionMutex.Unlock()
					continue
				}

				if !hopReturn.Success {
					continue
				}

				mtrMutex.Lock()
				mtrReturn := mtrReturns[ttl]
				mtrReturn.host = hopReturn.Addr
				mtrReturn.lastTime = hopReturn.Elapsed
				mtrReturn.allTime = append(mtrReturn.allTime, hopReturn.Elapsed)
				mtrReturn.succSum = mtrReturn.succSum + 1

				if mtrReturn.worstTime == time.Duration(0) || hopReturn.Elapsed > mtrReturn.worstTime {
					mtrReturn.worstTime = hopReturn.Elapsed
				}
				if mtrReturn.bestTime == time.Duration(0) || hopReturn.Elapsed < mtrReturn.bestTime {
					mtrReturn.bestTime = hopReturn.Elapsed
				}
				mtrReturn.sumTime += hopReturn.Elapsed
				mtrReturn.avgTime = time.Duration((int64)(mtrReturn.sumTime/time.Microsecond)/(int64)(mtrReturn.succSum)) * time.Microsecond
				mtrReturn.success = true
				mtrMutex.Unlock()

				if common.IsEqualIP(hopReturn.Addr, destAddr) {
					// could add early stop logic here
				}
			}
		}()
	}

	wg.Wait()

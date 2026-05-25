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

	// build results
	for index, mtrReturn := range mtrReturns {
		if index == 0 {
			continue
		}

		if mtrReturn == nil {
			break
		}

		hop := common.IcmpHop{TTL: mtrReturn.ttl, Snt: count}
		if index != 1 {
			hop.AddressFrom = mtrReturns[index-1].host
		} else {
			hop.AddressFrom = mtrReturn.host
		}
		hop.AddressTo = mtrReturn.host
		hop.Success = mtrReturn.success
		hop.LastTime = mtrReturn.lastTime
		hop.SumTime = mtrReturn.sumTime
		hop.AvgTime = mtrReturn.avgTime
		hop.BestTime = mtrReturn.bestTime
		hop.WorstTime = mtrReturn.worstTime

		if len(mtrReturn.allTime) > 0 {
			hop.SquaredDeviationTime = time.Duration(common.TimeSquaredDeviation(mtrReturn.allTime))
			hop.UncorrectedSDTime = time.Duration(common.TimeUncorrectedDeviation(mtrReturn.allTime))
			hop.CorrectedSDTime = time.Duration(common.TimeCorrectedDeviation(mtrReturn.allTime))
			hop.RangeTime = time.Duration(common.TimeRange(mtrReturn.allTime))
		}

		failSum := count - mtrReturn.succSum
		hop.SntFail = failSum
		loss := (float64)(failSum) / (float64)(count)
		hop.Loss = float64(loss)

		result.Hops = append(result.Hops, hop)

		if hop.Success {
			summaryKey := fmt.Sprintf("%d_%s", hop.TTL, hop.AddressTo)
			result.HopSummaryMap[summaryKey] = &common.IcmpSummary{
				AddressFrom: hop.AddressFrom,
				AddressTo:   hop.AddressTo,
				Snt:         hop.Snt,
				SntFail:     hop.SntFail,
				SntTime:     hop.SumTime,
			}
		}

		if common.IsEqualIP(hop.AddressTo, destAddr) {
			break
		}
	}

	return result
}

// RunMTRSequential execute sequential MTR route tracing
func RunMTRSequential(destAddr, srcAddr string, timeout time.Duration, maxHops, count int) *MtrResult {
	result := &MtrResult{
		DestAddr:      destAddr,
		Hops:          []common.IcmpHop{},
		HopSummaryMap: make(map[string]*common.IcmpSummary),
	}

	icmpID := &common.IcmpID{}
	pid := int(icmpID.Get())

	mtrReturns := make([]*MtrReturn, maxHops+1)

	seq := 0
	hasPermissionError := false

	for snt := 0; snt < count; snt++ {
		for ttl := 1; ttl < maxHops; ttl++ {
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
				}
			}

			hopReturn, err := icmp.Icmp(destAddr, srcAddr, ttl, pid, timeout, seq, false)
			seq++

			if err != nil {
				if !hasPermissionError {
					hasPermissionError = true
				}
				continue
			}

			if !hopReturn.Success {
				continue
			}

			mtrReturns[ttl].host = hopReturn.Addr
			mtrReturns[ttl].lastTime = hopReturn.Elapsed
			mtrReturns[ttl].allTime = append(mtrReturns[ttl].allTime, hopReturn.Elapsed)
			mtrReturns[ttl].succSum = mtrReturns[ttl].succSum + 1

			if mtrReturns[ttl].worstTime == time.Duration(0) || hopReturn.Elapsed > mtrReturns[ttl].worstTime {
				mtrReturns[ttl].worstTime = hopReturn.Elapsed
			}
			if mtrReturns[ttl].bestTime == time.Duration(0) || hopReturn.Elapsed < mtrReturns[ttl].bestTime {
				mtrReturns[ttl].bestTime = hopReturn.Elapsed
			}
			mtrReturns[ttl].sumTime += hopReturn.Elapsed
			mtrReturns[ttl].avgTime = time.Duration((int64)(mtrReturns[ttl].sumTime/time.Microsecond)/(int64)(mtrReturns[ttl].succSum)) * time.Microsecond
			mtrReturns[ttl].success = true

			if common.IsEqualIP(hopReturn.Addr, destAddr) {
				break
			}
		}
	}

	// build results (same as concurrent version)
	for index, mtrReturn := range mtrReturns {
		if index == 0 {
			continue
		}

		if mtrReturn == nil {
			break
		}

		hop := common.IcmpHop{TTL: mtrReturn.ttl, Snt: count}
		if index != 1 {
			hop.AddressFrom = mtrReturns[index-1].host
		} else {
			hop.AddressFrom = mtrReturn.host
		}
		hop.AddressTo = mtrReturn.host
		hop.Success = mtrReturn.success
		hop.LastTime = mtrReturn.lastTime
		hop.SumTime = mtrReturn.sumTime
		hop.AvgTime = mtrReturn.avgTime
		hop.BestTime = mtrReturn.bestTime
		hop.WorstTime = mtrReturn.worstTime

		if len(mtrReturn.allTime) > 0 {
			hop.SquaredDeviationTime = time.Duration(common.TimeSquaredDeviation(mtrReturn.allTime))
			hop.UncorrectedSDTime = time.Duration(common.TimeUncorrectedDeviation(mtrReturn.allTime))
			hop.CorrectedSDTime = time.Duration(common.TimeCorrectedDeviation(mtrReturn.allTime))
			hop.RangeTime = time.Duration(common.TimeRange(mtrReturn.allTime))
		}

		failSum := count - mtrReturn.succSum
		hop.SntFail = failSum
		loss := (float64)(failSum) / (float64)(count)
		hop.Loss = float64(loss)

		result.Hops = append(result.Hops, hop)

		if hop.Success {
			summaryKey := fmt.Sprintf("%d_%s", hop.TTL, hop.AddressTo)
			result.HopSummaryMap[summaryKey] = &common.IcmpSummary{
				AddressFrom: hop.AddressFrom,
				AddressTo:   hop.AddressTo,
				Snt:         hop.Snt,
				SntFail:     hop.SntFail,
				SntTime:     hop.SumTime,
			}
		}

		if common.IsEqualIP(hop.AddressTo, destAddr) {
			break
		}
	}

	return result
}

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

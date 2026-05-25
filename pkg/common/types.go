package common

import (
	"sync/atomic"
	"time"
)

type IcmpReturn struct {
	Success bool
	Addr    string
	Elapsed time.Duration
}

type IcmpID struct {
	icmpID int32
}

func (c *IcmpID) Get() int32 {
	for {
		val := atomic.LoadInt32(&c.icmpID)
		if val == 0 {
			atomic.StoreInt32(&c.icmpID, 1)
			val = 1
		}
		if atomic.CompareAndSwapInt32(&c.icmpID, 65500, 2) {
			return 1
		}
		if atomic.CompareAndSwapInt32(&c.icmpID, val, val+1) {
			return val
		}
	}
}

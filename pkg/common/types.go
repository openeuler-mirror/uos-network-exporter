package common

import (
	"sync/atomic"
)

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

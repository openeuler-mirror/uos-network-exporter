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

type IcmpSummary struct {
	AddressFrom string        `json:"address_from"`
	AddressTo   string        `json:"address_to"`
	Snt         int           `json:"snt"`
	SntFail     int           `json:"snt_fail"`
	SntTime     time.Duration `json:"snt_time"`
}

type IcmpHop struct {
	Success     bool          `json:"success"`
	AddressFrom string        `json:"address_from"`
	AddressTo   string        `json:"address_to"`
	TTL         int           `json:"ttl"`
	Snt         int           `json:"snt"`
	SntFail     int           `json:"snt_fail"`
	LastTime    time.Duration `json:"last"`
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

package mtr

import (
	"time"

	"uos_network_exporter/pkg/common"
)

const defaultMaxHops = 30
const defaultTimeout = 5 * time.Second
const defaultPackerSize = 56
const defaultCount = 10

// MtrResult Calculated results
type MtrResult struct {
	DestAddr      string                         `json:"dest_address"`
	Hops          []common.IcmpHop               `json:"hops"`
	HopSummaryMap map[string]*common.IcmpSummary `json:"hop_summary_map"`
}

// MtrReturn MTR Response
type MtrReturn struct {
	success   bool
	ttl       int
	host      string
	succSum   int
	lastTime  time.Duration
	allTime   []time.Duration
	sumTime   time.Duration
	bestTime  time.Duration
	avgTime   time.Duration
	worstTime time.Duration
}

// MtrOptions MTR Options
type MtrOptions struct {
	maxHops    int
	timeout    time.Duration
	packetSize int
	count      int
}

// MaxHops Getter
func (options *MtrOptions) MaxHops() int {
	if options.maxHops == 0 {
		options.maxHops = defaultMaxHops
	}
	return options.maxHops
}

// SetMaxHops Setter
func (options *MtrOptions) SetMaxHops(maxHops int) {
	options.maxHops = maxHops
}

// Timeout Getter
func (options *MtrOptions) Timeout() time.Duration {
	if options.timeout == 0 {
		options.timeout = defaultTimeout
	}
	return options.timeout
}

// SetTimeout Setter
func (options *MtrOptions) SetTimeout(timeout time.Duration) {
	options.timeout = timeout
}

// Count Getter
func (options *MtrOptions) Count() int {
	if options.count == 0 {
		options.count = defaultCount
	}
	return options.count
}

// SetCount Setter
func (options *MtrOptions) SetCount(count int) {
	options.count = count
}

// PacketSize Getter
func (options *MtrOptions) PacketSize() int {
	if options.packetSize == 0 {
		options.packetSize = defaultPackerSize
	}
	return options.packetSize
}

// SetPacketSize Setter
func (options *MtrOptions) SetPacketSize(packetSize int) {
	options.packetSize = packetSize
}

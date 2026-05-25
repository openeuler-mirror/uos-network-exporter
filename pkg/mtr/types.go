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

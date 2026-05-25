package tcp

import (
	"time"
)

// TestTCPPort simplified TCP port test implementation
func TestTCPPort(host, destIP, destPort, srcIP string, timeout time.Duration) *TCPPortReturn {
	// TODO: implement real TCP connection test logic
	result := &TCPPortReturn{
		Success:  true,
		DestAddr: host,
		DestIp:   destIP,
		DestPort: destPort,
		SrcIp:    srcIP,
		ConTime:  time.Millisecond * 50,
	}

	return result
}

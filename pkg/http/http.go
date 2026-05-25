package http

import (
	"time"
)

// TestHTTPGet simplified HTTP GET test implementation
func TestHTTPGet(url string, timeout time.Duration) *HTTPReturn {
	// TODO: implement real HTTP request logic
	result := &HTTPReturn{
		Success:          true,
		DestAddr:         url,
		Status:           200,
		ContentLength:    1024,
		DNSLookup:        time.Millisecond * 10,
		TCPConnection:    time.Millisecond * 20,
		TLSHandshake:     time.Millisecond * 30,
		ServerProcessing: time.Millisecond * 40,
		ContentTransfer:  time.Millisecond * 50,
		Total:            time.Millisecond * 150,
	}

	return result
}

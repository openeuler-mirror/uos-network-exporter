package common

import (
	"context"
	"fmt"
	"math"
	"net"
	"time"
)

func DestAddrs(ctx context.Context, host string, resolver *net.Resolver, timeout time.Duration) ([]string, error) {
	ipAddrs := make([]string, 0)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	addrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolving target: %v", err)
	}
	for _, addr := range addrs {
		ipAddr, err := net.ResolveIPAddr("ip", addr.IP.String())
		if err != nil {
			continue
		}
		ipAddrs = append(ipAddrs, ipAddr.IP.String())
	}
	return ipAddrs, nil
}

func IsEqualIP(ips1, ips2 string) bool {
	ip1 := net.ParseIP(ips1)
	if ip1 == nil {
		return false
	}
	ip2 := net.ParseIP(ips2)
	if ip2 == nil {
		return false
	}
	return ip1.String() == ip2.String()
}

// Time2Float Convert time to float32
func Time2Float(t time.Duration) float32 {
	return (float32)(t/time.Microsecond) / float32(1000)
}

// TimeRange finds the range of a slice of durations
func TimeRange(values []time.Duration) time.Duration {
	if len(values) <= 1 {
		return time.Duration(0)
	}
	min := values[0]
	max := time.Duration(0)
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

// TimeAverage Calculates the average of a slice of durations
func TimeAverage(values []time.Duration) float64 {
	l := len(values)
	if l <= 0 {
		return float64(0.0)
	}
	s := time.Duration(0)
	for _, d := range values {
		s += d
	}
	return float64(s) / float64(l)
}

// TimeSquaredDeviation Calculates the squared deviation
func TimeSquaredDeviation(values []time.Duration) float64 {
	avg := TimeAverage(values)
	sd := 0.0
	for _, v := range values {
		sd += math.Pow((float64(v) - float64(avg)), 2.0)
	}
	return sd
}

// TimeUncorrectedDeviation Calculates standard deviation without correction
func TimeUncorrectedDeviation(values []time.Duration) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sd := TimeSquaredDeviation(values)
	return math.Sqrt(sd / float64(len(values)))
}
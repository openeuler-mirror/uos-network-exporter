package icmp

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"

	"uos_network_exporter/pkg/common"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// https://hechao.li/2018/09/27/How-Is-Ping-Deduplexed/
const (
	protocolICMP     = 1  // Internet Control Message
	protocolIPv6ICMP = 58 // ICMP for IPv6
)

// Icmp execute real ICMP test
func Icmp(destAddr string, srcAddr string, ttl int, pid int, timeout time.Duration, seq int, ipv6 bool) (hop common.IcmpReturn, err error) {
	dstIp := net.ParseIP(destAddr)
	if dstIp == nil {
		return hop, fmt.Errorf("destination ip: %v is invalid", destAddr)
	}

	ipAddr := net.IPAddr{IP: dstIp}

	if srcAddr != "" {
		srcIp := net.ParseIP(srcAddr)
		if srcIp == nil {
			return hop, fmt.Errorf("source ip: %v is invalid, target: %v", srcAddr, destAddr)
		}

		if p4 := dstIp.To4(); len(p4) == net.IPv4len {
			return icmpIpv4(srcAddr, &ipAddr, ttl, pid, timeout, seq)
		}
		if ipv6 {
			return icmpIpv6(srcAddr, &ipAddr, ttl, pid, timeout, seq)
		} else {
			return hop, nil
		}
	}

	if p4 := dstIp.To4(); len(p4) == net.IPv4len {
		return icmpIpv4("0.0.0.0", &ipAddr, ttl, pid, timeout, seq)
	}
	if ipv6 {
		return icmpIpv6("::", &ipAddr, ttl, pid, timeout, seq)
	} else {
		return hop, nil
	}
}

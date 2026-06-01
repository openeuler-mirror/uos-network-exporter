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

func icmpIpv4(localAddr string, dst net.Addr, ttl int, pid int, timeout time.Duration, seq int) (hop common.IcmpReturn, err error) {
	hop.Success = false
	start := time.Now()
	c, err := icmp.ListenPacket("ip4:icmp", localAddr)
	if err != nil {
		return hop, err
	}
	defer c.Close()

	if err = c.IPv4PacketConn().SetTTL(ttl); err != nil {
		return hop, err
	}

	if err = c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return hop, err
	}

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, safeIntToUint32(seq))
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   pid,
			Seq:  seq,
			Data: append(bs, 'x'),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		return hop, err
	}

	if _, err := c.WriteTo(wb, dst); err != nil {
		return hop, err
	}

	peer, _, err := listenForSpecific4(c, append(bs, 'x'), pid, seq, wb)
	if err != nil {
		return hop, err
	}

	elapsed := time.Since(start)
	hop.Elapsed = elapsed
	hop.Addr = peer
	hop.Success = true
	return hop, err
}

package common

import (
	"context"
	"fmt"
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

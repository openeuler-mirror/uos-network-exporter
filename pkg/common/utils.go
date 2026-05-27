package common

import (
	"context"
	"fmt"
	"net"
	"time"
)

func DestAddrs(ctx context.Context, host string, resolver *net.Resolver, timeout time.Duration) ([]string, error) {

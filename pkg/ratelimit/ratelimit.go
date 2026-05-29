package ratelimit

import (
	"errors"
	"time"
)

var (
	ErrRateLimited   = errors.New("rate limited")
	ErrRateLimitSize = errors.New("limit must be greater than zero")
	ErrRateLimitTime = errors.New("invalid limit")
)

type RateLimiter struct {
	tokens chan struct{}
	limit  time.Duration
	ticker *time.Ticker
}

func NewRateLimiter(limit time.Duration, chanSize int) (*RateLimiter, error) {
	if chanSize <= 0 {
		return nil, ErrRateLimitSize
	}
	if limit <= 0 {
		return nil, ErrRateLimitTime
	}
	rl := &RateLimiter{
		tokens: make(chan struct{}, chanSize),
		limit:  limit,
		ticker: time.NewTicker(limit),
	}

	for i := 0; i < chanSize; i++ {
		rl.tokens <- struct{}{}
	}

	go rl.startRefreshTokens()
	return rl, nil
}

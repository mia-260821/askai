package utils

import (
	"context"
	"log"
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
	WaitWithTimeout(ctx context.Context, expiry time.Duration) (bool, error)
}

type TokenBucket struct {
	capacity int
	rate     int
	timeUnit time.Duration
	balance  int

	mutex sync.Mutex
}

func NewTokenBucket(ctx context.Context, capacity int, rate int, timeUnit time.Duration) RateLimiter {
	bucket := &TokenBucket{
		capacity: capacity,
		rate:     rate,
		timeUnit: timeUnit,
		balance:  0,
	}
	bucket.refill()
	go bucket.start(ctx)
	return bucket
}

func (t *TokenBucket) WaitWithTimeout(ctx context.Context, expiry time.Duration) (bool, error) {
	intervalTicker := time.NewTicker(t.timeUnit / 100)
	timeoutCtx, cancel := context.WithTimeout(ctx, expiry)
	defer cancel()
	for {
		select {
		case <-intervalTicker.C:
			if t.Allow() {
				return true, nil
			}
		case <-timeoutCtx.Done():
			if err := ctx.Err(); err != nil {
				return false, err
			}
			return false, nil
		}
	}
}

func (t *TokenBucket) Allow() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.balance > 0 {
		t.balance--
		return true
	}
	return false
}

func (t *TokenBucket) refill() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.balance+t.rate >= t.capacity {
		t.balance = t.capacity
	} else {
		t.balance += t.rate
	}
}

func (t *TokenBucket) start(ctx context.Context) {
	ticker := time.NewTicker(t.timeUnit)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				log.Println("context cancelled with an error", ctx.Err().Error())
			}
			return
		case <-ticker.C:
			t.refill()
		}
	}
}

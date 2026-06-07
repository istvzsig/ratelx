package ratelx

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu            sync.Mutex
	capacity      float64
	tokens        float64
	ratePerSecond float64
	lastRefill    time.Time
	debug         bool
}

func New(capacity int, rate float64, debug bool) *Limiter {
	return &Limiter{
		capacity:      float64(capacity),
		tokens:        float64(capacity),
		ratePerSecond: rate,
		lastRefill:    time.Now(),
		debug:         debug,
	}
}

// Allow checks if request is allowed immediately
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	return false
}

// Wait blocks until request is allowed
func (l *Limiter) Wait(ctx context.Context) error {
	timer := time.NewTimer(50 * time.Millisecond)

	delay := 10 * time.Millisecond
	if l.tokens < 0.5 {
		delay = 50 * time.Millisecond
	}

	defer timer.Stop()

	for {
		l.mu.Lock()
		l.refill()

		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}

		l.mu.Unlock()

		timer.Reset(delay)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}

// refill tokens over time
func (l *Limiter) refill() {
	now := time.Now()

	elapsed := now.Sub(l.lastRefill).Seconds()
	l.lastRefill = now

	l.tokens += elapsed * l.ratePerSecond

	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}

	if l.debug {
		fmt.Println("tokens:", l.tokens)
	}
}

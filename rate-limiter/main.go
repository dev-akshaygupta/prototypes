package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Simplest Implementation
type TokenBucket struct {
	capacity   int        // max tokens
	tokens     int        // current tokens
	refillRate int        // tokens added per second
	lastRefill time.Time  // last time token was refilled
	mu         sync.Mutex // protect bucket state
}

var wg sync.WaitGroup

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // starting full
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// refill tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsedTime := now.Sub(tb.lastRefill).Seconds()

	// how many tokens to add
	newTokens := int(elapsedTime * float64(tb.refillRate))
	if newTokens > 0 {
		tb.tokens += newTokens
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		// fmt.Printf("Refilling tokens %d\n", elapsedTime)
		tb.lastRefill = now
	}
}

// try to consume 1 token, Return true if allowed, else false.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// referesh bucket based on elapsed time
	tb.refill()

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

func main() {
	bucket := NewTokenBucket(5, 3) // Max Token = 5 and refill every 1 sec

	numReq := 30
	min := 50
	max := 200
	sleepMs := rand.Intn(max-min+1) + min
	for i := range numReq {
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		wg.Go(func() {
			if bucket.Allow() {
				fmt.Printf("Serving Request #%d at : %s\n", i, time.Now().Format("15:04:05.000"))
			} else {
				fmt.Printf("No tokens are available. Rejected Request #%d at : %s\n", i, time.Now().Format("15:04:05.000"))
			}
		})
	}
	wg.Wait()
}

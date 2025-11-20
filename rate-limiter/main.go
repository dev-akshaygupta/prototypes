package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var mu sync.Mutex

// 5 tokens = at most 5 concurrent requests
var tokens = map[string]bool{
	"token1": false,
	"token2": false,
	"token3": false,
	"token4": false,
	"token5": false,
}

var processed = 0

func APIServer(token string, reqNum int) {
	fmt.Printf("Request %d being processed with token: %s\n", reqNum, token)
	processed += 1

	// To mimic some proceess is happening
	time.Sleep(3 * time.Millisecond)
}

// In case of rate limiter tokens map is a shared resource
// APIServer is a request processor
func TokenRateLimiter(reqNum int) {

	// Acquire a token
	mu.Lock()
	var token string
	for t := range tokens {
		if !tokens[t] {
			token = t
			tokens[t] = true
			break
		}
	}

	// If no tokens are available: wait until one becomes free
	for token == "" {
		mu.Unlock() // release the lock for now

		fmt.Printf("Request %d is waiting for token...!\n", reqNum)
		time.Sleep(1 * time.Millisecond) // wait + backoff

		// Acquire a token
		mu.Lock()
		for t := range tokens {
			if !tokens[t] {
				token = t
				tokens[t] = true
				break
			}
		}
	}
	mu.Unlock()

	APIServer(token, reqNum)

	mu.Lock()
	tokens[token] = false
	mu.Unlock()
	// wg.Done()
}

func main() {
	numReq := 15

	for i := range numReq {
		wg.Go(func() {
			TokenRateLimiter(i)
		})
		// wg.Add(1)
		// go RateLimiter(i)
	}
	wg.Wait()

	defer fmt.Printf("\nTotal number of requests proceessed: %d\n", processed)
}

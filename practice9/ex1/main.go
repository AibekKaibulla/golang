package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	if rand.Intn(10) < 7 {
		fmt.Println("Operation failed, retrying...")
		return errors.New("temporary failure")
	}
	fmt.Println("Operation succeeded!")
	return nil
}

// func naiveInfiniteLoop() {
// 	var err error
// 	for {
// 		err = doSomethingUnreliable()
// 		if err == nil {
// 		break // Success, breaking the loop
// 		}
// 		// Error, just try again immediately
// 	}
// 	if err != nil {
// 	fmt.Printf("Failed after multiple retries: %v\n", err)
// 	}
// }

func fixedDelayRetry() {
	fmt.Println("\n=== Example 2: Fixed Delay Retry ===")
	var err error
	const maxRetries = 5
	const delay = 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d failed, waiting %v before next retry...\n", attempt+1, delay)
		if attempt < maxRetries-1 {
			time.Sleep(delay)
		}
	}
	if err != nil {
		fmt.Printf("Failed after %d attempts: %v\n", maxRetries, err)
	} else {
		fmt.Println("Succeeded within retry limit.")
	}
}

func exponentialBackoffRetry() {
	fmt.Println("\n=== Example 3: Exponential Backoff ===")
	var err error
	const maxRetries = 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			break
		}
		if attempt == maxRetries-1 {
			break
		}
		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}
		fmt.Printf("Attempt %d failed, waiting %v before next retry...\n", attempt+1, backoffTime)
		time.Sleep(backoffTime)
	}
	if err != nil {
		fmt.Printf("Failed after %d attempts: %v\n", maxRetries, err)
	} else {
		fmt.Println("Succeeded.")
	}
}

func exponentialBackoffWithJitter() {
	fmt.Println("\n=== Example 4: Exponential Backoff + Full Jitter (BEST) ===")
	var err error
	const maxRetries = 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			break
		}
		if attempt == maxRetries-1 {
			break
		}
		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoffTime)))
		fmt.Printf("Attempt %d failed, waiting ~%v (max backoff: %v + jitter) before next retry...\n",
			attempt+1, jitter, backoffTime)
		time.Sleep(jitter)
	}
	if err != nil {
		fmt.Printf("Failed after %d attempts: %v\n", maxRetries, err)
	} else {
		fmt.Println("Succeeded.")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// naiveInfiniteLoop()
	fixedDelayRetry()
	exponentialBackoffRetry()
	exponentialBackoffWithJitter()
}
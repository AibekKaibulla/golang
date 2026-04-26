package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"
)

type RetryClient struct {
	HTTPClient *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

func (c *RetryClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	for attempt := 1; attempt <= c.MaxRetries; attempt++ {
		reqClone := req.Clone(ctx)
		resp, err := c.HTTPClient.Do(reqClone)

		if err == nil && resp != nil && resp.StatusCode < http.StatusBadRequest {
			fmt.Printf("Attempt %d: Success!\n", attempt)
			return resp, nil
		}

		if !IsRetryable(resp, err) {
			return resp, err
		}

		if attempt == c.MaxRetries {
			if resp != nil && resp.Body != nil {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
			return nil, fmt.Errorf("all %d attempts failed", c.MaxRetries)
		}

		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		backoff := c.backoff(attempt)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt, backoff.Round(time.Millisecond))

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
	return nil, fmt.Errorf("max retries exceeded")
}

func (c *RetryClient) backoff(attempt int) time.Duration {
	exp := c.BaseDelay * time.Duration(1<<uint(attempt-1))
	const maxDelay = 5 * time.Second
	if exp > maxDelay {
		exp = maxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(exp / 2)))
	return exp + jitter
}

func main() {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&requestCount, 1)
		if n <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	client := &RetryClient{
		HTTPClient: &http.Client{},
		MaxRetries: 5,
		BaseDelay:  500 * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		fmt.Printf("failed to build request: %v\n", err)
		return
	}

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		fmt.Printf("request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("Server response: %v\n", result)
}

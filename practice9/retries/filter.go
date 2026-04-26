package main

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"time"
)


func IsRetryable(resp *http.Response, err error) bool {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		case http.StatusUnauthorized,
			http.StatusNotFound:
			return false
		}
	}

	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}

func CalculateBackoff(attempt int) time.Duration {
	var jts = rand.New(rand.NewSource(time.Now().UnixNano()))

	const (
		base_delay = 100 * time.Millisecond
		mx_delay  = 5 * time.Second
		mx_shift  = 7
	)

	if attempt > mx_shift {
		attempt = mx_shift
	}

	cap_delay := base_delay << uint(attempt)
	if cap_delay > mx_delay {
		cap_delay = mx_delay
	}

	return time.Duration(jts.Int63n(int64(cap_delay)))
}

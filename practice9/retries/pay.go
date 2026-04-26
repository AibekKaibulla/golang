package main

import (
	"context"
	"net/http"
	"errors"
	"io"
	"time"
)

func ExecutePayment(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		client = http.DefaultClient
	}
	if req == nil {
		return nil, errors.New("request is nil")
	}

	var attempt int

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		reqWithCtx := req.Clone(ctx)
		resp, err := client.Do(reqWithCtx)
		if err == nil && resp != nil && resp.StatusCode < http.StatusBadRequest {
			return resp, nil
		}

		if !IsRetryable(resp, err) {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		backoff := CalculateBackoff(attempt)
		attempt++


		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

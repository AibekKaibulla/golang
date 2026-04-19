package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetRate(t *testing.T) {
	t.Run("Successful scenario", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.85}`)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		rate, err := service.GetRate("USD", "EUR")
		assert.NoError(t, err)
		assert.Equal(t, 0.85, rate)
	})

	t.Run("API Business Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"invalid currency pair"}`)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "INVALID")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api error: invalid currency pair")
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `invalid json`)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode error")
	})

	t.Run("Slow Response/Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(6 * time.Second)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.85}`)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("Server Panic / 500 Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"internal server error"}`)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api error: internal server error")
	})

	t.Run("Empty Body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Empty body
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode error")
	})
}
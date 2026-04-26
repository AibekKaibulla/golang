package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)


type RequestStatus string

const (
	StatusProcessing RequestStatus = "processing"
	StatusCompleted  RequestStatus = "completed"
)

type CachedResponse struct {
	Status     RequestStatus
	StatusCode int
	Body       []byte
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]*CachedResponse)}
}

func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; exists {
		return false
	}
	m.data[key] = &CachedResponse{Status: StatusProcessing}
	return true
}

func (m *MemoryStore) Finish(key string, statusCode int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = &CachedResponse{
		Status:     StatusCompleted,
		StatusCode: statusCode,
		Body:       body,
	}
}


func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, `{"error":"Idempotency-Key header required"}`, http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			switch cached.Status {
			case StatusCompleted:
				fmt.Printf("[Middleware] Key %s: returning cached result (business logic NOT re-executed)\n", key)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Idempotent-Replayed", "true")
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			case StatusProcessing:
				fmt.Printf("[Middleware] Key %s: still processing -> 409 Conflict\n", key)
				http.Error(w, `{"error":"request is already in progress"}`, http.StatusConflict)
				return
			}
		}

		if !store.StartProcessing(key) {
			if cached, exists := store.Get(key); exists && cached.Status == StatusCompleted {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Idempotent-Replayed", "true")
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				fmt.Printf("[Middleware] Key %s: race lost -> 409 Conflict\n", key)
				http.Error(w, `{"error":"request is already in progress"}`, http.StatusConflict)
			}
			return
		}

		fmt.Printf("[Middleware] Key %s: new request, processing started\n", key)

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		store.Finish(key, recorder.Code, recorder.Body.Bytes())
		fmt.Printf("[Middleware] Key %s: processing finished, result saved\n", key)

		for k, vals := range recorder.Header() {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}


type PaymentResponse struct {
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
	TransactionID string `json:"transaction_id"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Handler] Processing started (heavy operation)...")
	time.Sleep(2 * time.Second)

	txID := uuid.New().String()
	resp := PaymentResponse{
		Status:        "paid",
		Amount:        1000,
		TransactionID: txID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	fmt.Printf("[Handler] Done. transaction_id=%s\n", txID)
}


func sendRequest(client *http.Client, url, idempotencyKey, label string, wg *sync.WaitGroup) {
	defer wg.Done()

	req, _ := http.NewRequest("POST", url+"/pay", strings.NewReader(`{"amount":1000}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[%s] Error: %v\n", label, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	replayed := resp.Header.Get("X-Idempotent-Replayed")
	fmt.Printf("[%s] Status=%d replayed=%s body=%s\n", label, resp.StatusCode, replayed, strings.TrimSpace(string(body)))
}

func main() {
	store := NewMemoryStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/pay", paymentHandler)

	handler := IdempotencyMiddleware(store, mux)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fmt.Println("=== Task 2: Loan Repayment - Idempotency Middleware ===")
	fmt.Printf("Test server: %s\n\n", ts.URL)

	idempotencyKey := uuid.New().String()
	fmt.Printf("Using Idempotency-Key: %s\n\n", idempotencyKey)

	client := &http.Client{Timeout: 10 * time.Second}

	fmt.Println("--- Phase 1: Simulating double-click attack (7 concurrent requests) ---")
	var wg sync.WaitGroup
	for i := 1; i <= 7; i++ {
		wg.Add(1)
		go sendRequest(client, ts.URL, idempotencyKey, fmt.Sprintf("goroutine-%d", i), &wg)
	}
	wg.Wait()

	fmt.Println("\n--- Phase 2: Replay after completion (same key, business logic must NOT run again) ---")
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go sendRequest(client, ts.URL, idempotencyKey, "replay-after-done", &wg2)
	wg2.Wait()

	fmt.Println("\n--- Phase 3: Missing Idempotency-Key header -> 400 ---")
	req, _ := http.NewRequest("POST", ts.URL+"/pay", nil)
	resp, _ := client.Do(req)
	fmt.Printf("[no-key] Status=%d\n", resp.StatusCode)
}
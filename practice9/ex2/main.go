package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
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
	resp, exists := m.data[key]
	return resp, exists
}

func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; exists {
		return false
	}
	m.data[key] = &CachedResponse{Completed: false}
	return true
}

func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	if resp, exists := m.data[key]; exists {
		resp.StatusCode = status
		resp.Body = body
		resp.Completed = true
	} else {
		m.data[key] = &CachedResponse{StatusCode: status, Body: body, Completed: true}
	}
	m.mu.Unlock()
}

func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			}
			return
		}

		if !store.StartProcessing(key) {
			if cached, exists := store.Get(key); exists && cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			}
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		store.Finish(key, recorder.Code, recorder.Body.Bytes())

		for k, vals := range recorder.Header() {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}

func businessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"status": "ok", "message": "resource created"}`)
}

func main() {
	store := NewMemoryStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/resource", businessHandler)

	handler := IdempotencyMiddleware(store, mux)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fmt.Println("=== Part 2 Example: Idempotency Middleware ===")
	fmt.Printf("Test server: %s\n\n", ts.URL)

	idempotencyKey := "test-key-abc-123"

	fmt.Println("--- Request 1 (new key, should succeed) ---")
	req1, _ := http.NewRequest("POST", ts.URL+"/resource", strings.NewReader(`{"data":"value"}`))
	req1.Header.Set("Idempotency-Key", idempotencyKey)
	resp1, _ := http.DefaultClient.Do(req1)
	fmt.Printf("Status: %d\n\n", resp1.StatusCode)

	fmt.Println("--- Request 2 (same key, should return cached result) ---")
	req2, _ := http.NewRequest("POST", ts.URL+"/resource", strings.NewReader(`{"data":"value"}`))
	req2.Header.Set("Idempotency-Key", idempotencyKey)
	resp2, _ := http.DefaultClient.Do(req2)
	fmt.Printf("Status: %d (cached, business logic NOT re-executed)\n\n", resp2.StatusCode)

	fmt.Println("--- Request 3 (missing key, should get 400) ---")
	req3, _ := http.NewRequest("POST", ts.URL+"/resource", nil)
	resp3, _ := http.DefaultClient.Do(req3)
	fmt.Printf("Status: %d\n", resp3.StatusCode)
}
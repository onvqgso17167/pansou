package api

import (
	"testing"
	"time"
)

func TestSearchCache_SetAndGet(t *testing.T) {
	cache := NewSearchCache(5 * time.Minute)

	results := []map[string]interface{}{
		{"name": "file1", "url": "https://example.com/1"},
	}

	cache.Set("keyword:test", results)

	got, ok := cache.Get("keyword:test")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0]["name"] != "file1" {
		t.Errorf("unexpected result name: %v", got[0]["name"])
	}
}

func TestSearchCache_Miss(t *testing.T) {
	cache := NewSearchCache(5 * time.Minute)

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestSearchCache_Expiry(t *testing.T) {
	// Use a longer sleep (200ms) to reduce flaky failures on slow or busy CI machines.
	// The TTL is set to 50ms, so 200ms gives a comfortable margin.
	cache := NewSearchCache(50 * time.Millisecond)

	cache.Set("keyword:expire", []map[string]interface{}{{"name": "old"}})

	time.Sleep(200 * time.Millisecond)

	_, ok := cache.Get("keyword:expire")
	if ok {
		t.Fatal("expected expired cache entry to be a miss")
	}
}

func TestSearchCache_Delete(t *testing.T) {
	cache := NewSearchCache(5 * time.Minute)
	cache.Set("key", []map[string]interface{}{{"name": "item"}})
	cache.Delete("key")

	_, ok := cache.Get("key")
	if ok {
		t.Fatal("expected deleted key to be a miss")
	}
}

func TestSearchCache_Flush(t *testing.T) {
	cache := NewSearchCache(50 * time.Millisecond)

	cache.Set("k1", []map[string]interface{}{{"name": "a"}})
	cache.Set("k2", []map[string]interface{}{{"name": "b"}})

	// Wait for k1 and k2 to expire before adding k3
	time.Sleep(200 * time.Millisecond)

	cache.Set("k3", []map[string]interface{}{{"name": "c"}})

	cache.Flush()

	if cache.Len() != 1 {
		t.Fatalf("expected 1 entry after flush, got %d", cache.Len())
	}
}

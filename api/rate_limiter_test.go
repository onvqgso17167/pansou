package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestIPRateLimiter_Allow(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(2), 2)

	// First two requests should be allowed (burst=2).
	if !limiter.Allow("127.0.0.1") {
		t.Error("expected first request to be allowed")
	}
	if !limiter.Allow("127.0.0.1") {
		t.Error("expected second request to be allowed")
	}
	// Third request should be denied (burst exhausted).
	if limiter.Allow("127.0.0.1") {
		t.Error("expected third request to be denied")
	}
}

func TestIPRateLimiter_SeparateIPs(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(1), 1)

	if !limiter.Allow("10.0.0.1") {
		t.Error("expected 10.0.0.1 first request to be allowed")
	}
	// Different IP should have its own limiter.
	if !limiter.Allow("10.0.0.2") {
		t.Error("expected 10.0.0.2 first request to be allowed")
	}
}

func TestRateLimitMiddleware_Blocks(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(1), 1)
	mw := RateLimitMiddleware(limiter)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")

	// First request — allowed.
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	// Second request — rate limited.
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rr2.Code)
	}
}

func TestRealIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "1.2.3.4")
	if ip := realIP(req); ip != "1.2.3.4" {
		t.Errorf("expected 1.2.3.4, got %s", ip)
	}
}

func TestRealIP_Fallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "5.6.7.8:9000"
	if ip := realIP(req); ip != "5.6.7.8:9000" {
		t.Errorf("expected 5.6.7.8:9000, got %s", ip)
	}
}

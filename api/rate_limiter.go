package api

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter holds per-IP rate limiters.
type IPRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

// NewIPRateLimiter creates a new IPRateLimiter with the given rate and burst.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter returns the rate limiter for the given IP, creating one if needed.
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	if lim, exists := i.limiters[ip]; exists {
		return lim
	}

	lim := rate.NewLimiter(i.r, i.b)
	i.limiters[ip] = lim
	return lim
}

// Allow reports whether the given IP is allowed to make a request.
func (i *IPRateLimiter) Allow(ip string) bool {
	return i.getLimiter(ip).Allow()
}

// Cleanup removes stale limiters (for long-running servers).
func (i *IPRateLimiter) Cleanup(ttl time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()
	// Simple reset — replace with time-tracked map for production use.
	i.limiters = make(map[string]*rate.Limiter)
}

// RateLimitMiddleware returns an HTTP middleware that enforces per-IP rate limits.
func RateLimitMiddleware(limiter *IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := realIP(r)
			if !limiter.Allow(ip) {
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// realIP extracts the real client IP from common headers or RemoteAddr.
// When X-Forwarded-For contains a chain of IPs, use only the first (client) IP.
func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For may be "client, proxy1, proxy2" — take the first entry
		return strings.TrimSpace(strings.SplitN(ip, ",", 2)[0])
	}
	return r.RemoteAddr
}

package api

import (
	"os"
	"strconv"

	"golang.org/x/time/rate"
)

// RateLimiterConfig holds configuration for the rate limiter.
type RateLimiterConfig struct {
	// RequestsPerSecond is the sustained request rate allowed per IP.
	RequestsPerSecond float64
	// Burst is the maximum burst size allowed per IP.
	Burst int
}

// DefaultRateLimiterConfig returns sensible defaults.
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 5,
		Burst:             10,
	}
}

// RateLimiterConfigFromEnv reads rate limiter settings from environment variables:
//   RATE_LIMIT_RPS  — requests per second (float, default 5)
//   RATE_LIMIT_BURST — burst size (int, default 10)
func RateLimiterConfigFromEnv() RateLimiterConfig {
	cfg := DefaultRateLimiterConfig()

	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if v, err := strconv.ParseFloat(rps, 64); err == nil && v > 0 {
			cfg.RequestsPerSecond = v
		}
	}

	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		if v, err := strconv.Atoi(burst); err == nil && v > 0 {
			cfg.Burst = v
		}
	}

	return cfg
}

// BuildLimiter constructs an IPRateLimiter from the config.
func (c RateLimiterConfig) BuildLimiter() *IPRateLimiter {
	return NewIPRateLimiter(rate.Limit(c.RequestsPerSecond), c.Burst)
}

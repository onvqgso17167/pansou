package api

import (
	"os"
	"testing"
)

func TestDefaultRateLimiterConfig(t *testing.T) {
	cfg := DefaultRateLimiterConfig()
	if cfg.RequestsPerSecond != 5 {
		t.Errorf("expected RPS=5, got %f", cfg.RequestsPerSecond)
	}
	if cfg.Burst != 10 {
		t.Errorf("expected Burst=10, got %d", cfg.Burst)
	}
}

func TestRateLimiterConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("RATE_LIMIT_RPS")
	os.Unsetenv("RATE_LIMIT_BURST")

	cfg := RateLimiterConfigFromEnv()
	if cfg.RequestsPerSecond != 5 {
		t.Errorf("expected default RPS=5, got %f", cfg.RequestsPerSecond)
	}
	if cfg.Burst != 10 {
		t.Errorf("expected default Burst=10, got %d", cfg.Burst)
	}
}

func TestRateLimiterConfigFromEnv_Custom(t *testing.T) {
	os.Setenv("RATE_LIMIT_RPS", "20")
	os.Setenv("RATE_LIMIT_BURST", "50")
	defer os.Unsetenv("RATE_LIMIT_RPS")
	defer os.Unsetenv("RATE_LIMIT_BURST")

	cfg := RateLimiterConfigFromEnv()
	if cfg.RequestsPerSecond != 20 {
		t.Errorf("expected RPS=20, got %f", cfg.RequestsPerSecond)
	}
	if cfg.Burst != 50 {
		t.Errorf("expected Burst=50, got %d", cfg.Burst)
	}
}

func TestRateLimiterConfigFromEnv_Invalid(t *testing.T) {
	os.Setenv("RATE_LIMIT_RPS", "not-a-number")
	os.Setenv("RATE_LIMIT_BURST", "-5")
	defer os.Unsetenv("RATE_LIMIT_RPS")
	defer os.Unsetenv("RATE_LIMIT_BURST")

	cfg := RateLimiterConfigFromEnv()
	// Should fall back to defaults on invalid values.
	if cfg.RequestsPerSecond != 5 {
		t.Errorf("expected fallback RPS=5, got %f", cfg.RequestsPerSecond)
	}
	if cfg.Burst != 10 {
		t.Errorf("expected fallback Burst=10, got %d", cfg.Burst)
	}
}

func TestBuildLimiter(t *testing.T) {
	cfg := RateLimiterConfig{RequestsPerSecond: 3, Burst: 3}
	limiter := cfg.BuildLimiter()
	if limiter == nil {
		t.Fatal("expected non-nil limiter")
	}
	// Verify burst by consuming tokens.
	allowed := 0
	for i := 0; i < 5; i++ {
		if limiter.Allow("test-ip") {
			allowed++
		}
	}
	if allowed != 3 {
		t.Errorf("expected 3 allowed requests (burst), got %d", allowed)
	}
}

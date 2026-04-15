package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"pansou/api"
)

func main() {
	// Set Gin mode based on environment
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gint}
	gin.SetMode(mode)
 := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Apply rate limiting middleware
	rateCfg := api.RateLimiterConfigFromEnv()
	limiter := api.NewIPRateLimiter(rateCfg.Rate, rateCfg.Burst)
	r.Use(api.RateLimitMiddleware(limiter))

	// Initialize search cache
	cacheTTL := getCacheTTL()
	cache := api.NewSearchCache(cacheTTL)

	// Public routes
	v1 := r.Group("/api")
	{
		v1.GET("/search", api.SearchHandler(cache))
		v1.POST("/login", api.LoginHandler)
		v1.GET("/verify", api.VerifyHandler)
		v1.POST("/logout", api.LogoutHandler)
	}

	// Plugin management routes (protected)
	admin := r.Group("/admin")
	admin.Use(api.AuthMiddleware())
	{
		admin.GET("/plugins", api.ListPluginsHandler)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	port := getPort()
	addr := fmt.Sprintf(":%s", port)
	log.Printf("[pansou] starting server on %s (mode=%s)", addr, mode)

	if err := r.Run(addr); err != nil {
		log.Fatalf("[pansou] server error: %v", err)
	}
}

// getPort returns the HTTP port from the PORT environment variable,
// defaulting to 8080 if not set.
func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// getCacheTTL reads CACHE_TTL_SECONDS from the environment and returns
// the corresponding duration. Defaults to 5 minutes.
func getCacheTTL() time.Duration {
	if v := os.Getenv("CACHE_TTL_SECONDS"); v != "" {
		var secs int
		if _, err := fmt.Sscanf(v, "%d", &secs); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
		log.Printf("[pansou] invalid CACHE_TTL_SECONDS %q, using default", v)
	}
	return 5 * time.Minute
}

// AuthMiddleware returns a Gin middleware that validates the JWT token
// present in the Authorization header. It delegates to api.VerifyHandler
// logic but short-circuits the request on failure.
func AuthMiddleware() gin.HandlerFunc {
	return api.AuthMiddleware()
}
